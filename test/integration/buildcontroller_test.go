package integration

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kapi "k8s.io/kubernetes/pkg/api"
	kclientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"

	"github.com/openshift/origin/pkg/client"
	"github.com/openshift/origin/pkg/cmd/server/bootstrappolicy"
	"github.com/openshift/origin/pkg/cmd/server/origin"
	"github.com/openshift/origin/test/common/build"
	testutil "github.com/openshift/origin/test/util"
	testserver "github.com/openshift/origin/test/util/server"
)

type controllerCount struct {
	BuildControllers,
	BuildPodControllers,
	ImageChangeControllers,
	ConfigChangeControllers int
}

// TestConcurrentBuildControllers tests the transition of a build from new to pending. Ensures that only a single New -> Pending
// transition happens and that only a single pod is created during a set period of time.
func TestConcurrentBuildControllers(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	// Start a master with multiple BuildControllers
	osClient, kClient := setupBuildControllerTest(controllerCount{BuildControllers: 5}, t)
	build.RunBuildControllerTest(t, osClient, kClient)
}

// TestConcurrentBuildPodControllers tests the lifecycle of a build pod when running multiple controllers.
func TestConcurrentBuildPodControllers(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	// Start a master with multiple BuildPodControllers
	osClient, kClient := setupBuildControllerTest(controllerCount{BuildPodControllers: 5}, t)
	build.RunBuildPodControllerTest(t, osClient, kClient)
}

func TestConcurrentBuildImageChangeTriggerControllers(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	// Start a master with multiple ImageChangeTrigger controllers
	osClient, _ := setupBuildControllerTest(controllerCount{ImageChangeControllers: 5}, t)
	build.RunImageChangeTriggerTest(t, osClient)
}

func TestBuildDeleteController(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	osClient, kClient := setupBuildControllerTest(controllerCount{}, t)
	build.RunBuildDeleteTest(t, osClient, kClient)
}

func TestBuildRunningPodDeleteController(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	osClient, kClient := setupBuildControllerTest(controllerCount{}, t)
	build.RunBuildRunningPodDeleteTest(t, osClient, kClient)
}

func TestBuildCompletePodDeleteController(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	osClient, kClient := setupBuildControllerTest(controllerCount{}, t)
	build.RunBuildCompletePodDeleteTest(t, osClient, kClient)
}

func TestConcurrentBuildConfigControllers(t *testing.T) {
	defer testutil.DumpEtcdOnFailure(t)
	osClient, kClient := setupBuildControllerTest(controllerCount{ConfigChangeControllers: 5}, t)
	build.RunBuildConfigChangeControllerTest(t, osClient, kClient)
}

func setupBuildControllerTest(counts controllerCount, t *testing.T) (*client.Client, kclientset.Interface) {
	testutil.RequireEtcd(t)
	master, clusterAdminKubeConfig, err := testserver.StartTestMaster()
	if err != nil {
		t.Fatal(err)
	}

	clusterAdminClient, err := testutil.GetClusterAdminClient(clusterAdminKubeConfig)
	if err != nil {
		t.Fatal(err)
	}

	clusterAdminKubeClientset, err := testutil.GetClusterAdminKubeClient(clusterAdminKubeConfig)
	if err != nil {
		t.Fatal(err)
	}
	_, err = clusterAdminKubeClientset.Core().Namespaces().Create(&kapi.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: testutil.Namespace()},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testserver.WaitForServiceAccounts(clusterAdminKubeClientset, testutil.Namespace(), []string{bootstrappolicy.BuilderServiceAccountName, bootstrappolicy.DefaultServiceAccountName}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	openshiftConfig, err := origin.BuildMasterConfig(*master)
	if err != nil {
		t.Fatal(err)
	}

	// Get the build controller clients, since those rely on service account tokens
	// We don't want to proceed with the rest of the test until those are available
	openshiftConfig.BuildControllerClients()

	for i := 0; i < counts.BuildControllers; i++ {
		openshiftConfig.RunBuildController(openshiftConfig.Informers)
	}
	for i := 0; i < counts.BuildPodControllers; i++ {
		openshiftConfig.RunBuildPodController()
	}
	for i := 0; i < counts.ImageChangeControllers; i++ {
		openshiftConfig.RunImageTriggerController()
	}
	for i := 0; i < counts.ConfigChangeControllers; i++ {
		openshiftConfig.RunBuildConfigChangeController()
	}
	return clusterAdminClient, clusterAdminKubeClientset
}

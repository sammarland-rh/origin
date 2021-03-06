package cache

import (
	"k8s.io/client-go/tools/cache"

	templateapi "github.com/openshift/origin/pkg/template/api"
)

type StoreToTemplateLister interface {
	List() ([]*templateapi.Template, error)
	ListByNamespace(namespace string) ([]*templateapi.Template, error)
	GetByUID(uid string) (*templateapi.Template, error)
}

type StoreToTemplateListerImpl struct {
	cache.Indexer
}

func (s *StoreToTemplateListerImpl) List() ([]*templateapi.Template, error) {
	list := s.Indexer.List()

	templates := make([]*templateapi.Template, len(list))
	for i, template := range list {
		templates[i] = template.(*templateapi.Template)
	}
	return templates, nil
}

func (s *StoreToTemplateListerImpl) GetByUID(uid string) (*templateapi.Template, error) {
	templates, err := s.Indexer.ByIndex(TemplateUIDIndex, uid)
	if err != nil || len(templates) == 0 {
		return nil, err
	}
	return templates[0].(*templateapi.Template), nil
}

func (s *StoreToTemplateListerImpl) ListByNamespace(namespace string) ([]*templateapi.Template, error) {
	list, err := s.Indexer.ByIndex(cache.NamespaceIndex, namespace)
	if err != nil {
		return nil, err
	}

	templates := make([]*templateapi.Template, len(list))
	for i, template := range list {
		templates[i] = template.(*templateapi.Template)
	}

	return templates, nil
}

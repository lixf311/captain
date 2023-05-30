package horizontalpodautoscaler

import (
	"captain/pkg/bussiness/kube-resources/alpha1"
	"captain/pkg/unify/query"
	"captain/pkg/unify/response"
	autoV2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
)

type hpaV2beta2Provider struct {
	sharedInformers informers.SharedInformerFactory
}

func New(informer informers.SharedInformerFactory) hpaV2beta2Provider {
	return hpaV2beta2Provider{sharedInformers: informer}
}

func (cm hpaV2beta2Provider) Get(namespace, name string) (runtime.Object, error) {
	return cm.sharedInformers.Autoscaling().V2beta2().HorizontalPodAutoscalers().Lister().HorizontalPodAutoscalers(namespace).Get(name)
}

func (cm hpaV2beta2Provider) List(namespace string, query *query.QueryInfo) (*response.ListResult, error) {
	raw, err := cm.sharedInformers.Autoscaling().V2beta2().HorizontalPodAutoscalers().Lister().HorizontalPodAutoscalers(namespace).List(query.GetSelector())

	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, configMap := range raw {
		result = append(result, configMap)
	}

	return alpha1.DefaultList(result, query, compareFunc, filter), nil
}

func filter(object runtime.Object, filter query.Filter) bool {
	hpa, ok := object.(*autoV2.HorizontalPodAutoscaler)
	if !ok {
		return false
	}

	return alpha1.DefaultObjectMetaFilter(hpa.ObjectMeta, filter)
}

func compareFunc(left, right runtime.Object, field query.Field) bool {

	leftHPA, ok := left.(*autoV2.HorizontalPodAutoscaler)
	if !ok {
		return false
	}

	rightHPA, ok := right.(*autoV2.HorizontalPodAutoscaler)
	if !ok {
		return false
	}

	return alpha1.DefaultObjectMetaCompare(leftHPA.ObjectMeta, rightHPA.ObjectMeta, field)
}

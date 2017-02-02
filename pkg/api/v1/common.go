package v1

import (
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func RegisterType(clientset *kubernetes.Clientset, resourceName string) {
	log.Printf("Trying to register %s type\n", resourceName)

	// Check if resource is already registered.
	_, err := clientset.ExtensionsV1beta1().ThirdPartyResources().Get(resourceName, v1.GetOptions{})

	if err != nil {
		log.Printf("Registering new %s type\n", resourceName)
		tpr := &v1beta1.ThirdPartyResource{
			ObjectMeta: v1.ObjectMeta{
				Name: resourceName,
			},
			Versions: []v1beta1.APIVersion{
				{Name: APIVersion},
			},
			Description: "A specification of " + resourceName,
		}

		_, err := clientset.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
		if err != nil {
			log.Printf("Cannot register %s type\n", resourceName)
			panic(err.Error())
		}

		log.Printf("New %s type registered\n", resourceName)
	}
}

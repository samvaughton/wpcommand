package flow

import (
	"context"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func DeleteExpiredBuildPreviews(flowOpts types.FlowOptions) {
	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "DELETE_EXPIRED_BUILD_PREVIEWS",
		"Detail": "",
	}).Debug("started")

	client, err := kubernetes.NewForConfig(config.Config.K8RestConfig)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "DELETE_EXPIRED_BUILD_PREVIEWS",
			"Detail": "INIT_CONFIG",
		}).Error(err)
	}

	labelSelector := "wpcmd.k8.rentivo.com/type=preview"

	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "DELETE_EXPIRED_BUILD_PREVIEWS",
			"Detail": "LIST_NAMESPACES",
		}).Error(err)
	}

	for _, item := range namespaces.Items {
		log.WithFields(log.Fields{
			"Source":    flowOpts.LogSource,
			"Action":    "DELETE_EXPIRED_BUILD_PREVIEWS",
			"Detail":    "LIST_NAMESPACES",
			"Namespace": item.Name,
			"CreatedAt": item.CreationTimestamp.String(),
		}).Debug(fmt.Sprintf("found namespace"))

		// if creation + X hours STILL in the past, then we can delete as it is over X hours old
		if item.CreationTimestamp.Add(time.Hour * 4).Before(time.Now()) {
			err = client.CoreV1().Namespaces().Delete(context.Background(), item.Name, metav1.DeleteOptions{})

			log.WithFields(log.Fields{
				"Source":    flowOpts.LogSource,
				"Action":    "DELETE_EXPIRED_BUILD_PREVIEWS",
				"Detail":    "NAMESPACE_EXPIRY_CHECK",
				"Namespace": item.Name,
				"CreatedAt": item.CreationTimestamp.String(),
			}).Info(fmt.Sprintf("deleting namespace as it matches expiry criteria"))

			if err != nil {
				log.WithFields(log.Fields{
					"Source":    flowOpts.LogSource,
					"Action":    "DELETE_EXPIRED_BUILD_PREVIEWS",
					"Detail":    "NAMESPACE_DELETE",
					"Namespace": item.Name,
					"CreatedAt": item.CreationTimestamp.String(),
				}).Error(err)
			}
		}
	}

	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "DELETE_EXPIRED_BUILD_PREVIEWS",
		"Detail": "",
	}).Debug("finished")
}

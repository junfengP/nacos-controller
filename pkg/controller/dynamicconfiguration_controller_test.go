package controller

import (
	"fmt"
	v12 "github.com/nacos-group/nacos-controller/api/v1"
	"github.com/nacos-group/nacos-controller/pkg/nacos/auth"
	"github.com/nacos-group/nacos-controller/pkg/nacos/client/impl"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"math/rand"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	nacosServerAk        = "NACOS_SERVER_ACCESSKEY"
	nacosServerSk        = "NACOS_SERVER_SECRETKEY"
	nacosServerEndpoint  = "NACOS_SERVER_ENDPOINT"
	nacosServerNamespace = "NACOS_SERVER_NAMESPACE"
)

const (
	dcTestNamespaceStr        = "dc-suite-test"
	dcTestNacosCredentialName = "dc-suite-test"
	dcTestConfigMapName       = "dc-suite-test"
)

var _ = Describe("DynamicConfigurationController", func() {
	BeforeEach(func() {
		ensureEnvKeyExist(nacosServerAk)
		ensureEnvKeyExist(nacosServerSk)
		ensureEnvKeyExist(nacosServerEndpoint)
		ensureEnvKeyExist(nacosServerNamespace)
		ensureNamespaceExist(dcTestNamespaceStr)
		nacosCredential := v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dcTestNacosCredentialName,
				Namespace: dcTestNamespaceStr,
			},
			Data: map[string][]byte{
				"ak": []byte(os.Getenv(nacosServerAk)),
				"sk": []byte(os.Getenv(nacosServerSk)),
			},
		}
		ensureSecretExist(&nacosCredential)
	})
	Describe("DynamicConfiguration with Cluster2Server functional test", func() {
		It("SyncDeletion & SyncPolicyAlways", func() {
			rand.Seed(time.Now().Unix())
			randInt := rand.Int()
			dataId := fmt.Sprintf("randon-data-id-%d", randInt)
			content := fmt.Sprintf("content-%d", randInt)
			cm := v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      dcTestConfigMapName,
					Namespace: dcTestNamespaceStr,
				},
				Data: map[string]string{
					dataId: content,
				},
			}
			dc := v12.DynamicConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dc-suite-test-cluster2server",
					Namespace: dcTestNamespaceStr,
				},
				Spec: v12.DynamicConfigurationSpec{
					DataIds: []string{dataId},
					Strategy: v12.SyncStrategy{
						SyncPolicy:    v12.Always,
						SyncDeletion:  true,
						SyncDirection: v12.Cluster2Server,
					},
					NacosServer: v12.NacosServerConfiguration{
						Endpoint:  pointer.String(os.Getenv(nacosServerEndpoint)),
						Namespace: os.Getenv(nacosServerNamespace),
						Group:     fmt.Sprintf("suite-test-group-cluster2server-%d", randInt),
						AuthRef: &v1.ObjectReference{
							Name:       dcTestNacosCredentialName,
							APIVersion: "v1",
							Kind:       "Secret",
						},
					},
					ObjectRef: &v1.ObjectReference{
						Name:       dcTestConfigMapName,
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
				},
			}
			By("create a DynamicConfiguration with cluster2server sync direction")
			ensureConfigMapExist(&cm)
			ensureDynamicConfigurationExist(&dc)
			checkDynamicConfigurationStatus(&dc)

			By(fmt.Sprintf("check nacos server for dataId: %s", dataId))
			contentFromNacos, exist := getContentByDataId(&dc)
			gomega.Expect(exist).Should(gomega.BeTrue())
			gomega.Expect(contentFromNacos).Should(gomega.BeEquivalentTo(content), fmt.Sprintf("content mismatch, cluster: %s, nacos: %s", content, contentFromNacos))

			newRandInt := rand.Int()
			newContent := fmt.Sprintf("content-%d", newRandInt)
			By(fmt.Sprintf("modify configmap, change content to %s", newContent))
			cm.Data[dataId] = newContent
			ensureConfigMapExist(&cm)

			By(fmt.Sprintf("check nacos server for dataId: %s", dataId))
			gomega.Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      dc.Name,
					Namespace: dc.Namespace,
				}, &dc)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				checkDynamicConfigurationStatus(&dc)
				contentFromNacos, exist := getContentByDataId(&dc)
				gomega.Expect(exist).Should(gomega.BeTrue())
				return contentFromNacos == newContent
			}, time.Second*60, time.Second*5).Should(gomega.BeTrue())

			By("delete DynamicConfiguration")
			gomega.Expect(k8sClient.Delete(ctx, &dc)).Should(gomega.Succeed())
			gomega.Eventually(func() bool {
				v := v12.DynamicConfiguration{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: dc.Name, Namespace: dc.Namespace}, &v)
				if err == nil {
					return false
				}
				gomega.Expect(errors.IsNotFound(err)).Should(gomega.BeTrue())
				return true
			}, time.Second*30, time.Second*5).Should(gomega.BeTrue())

			By(fmt.Sprintf("check nacos server for dataId: %s", dataId))
			gomega.Eventually(func() bool {
				_, exist := getContentByDataId(&dc)
				return !exist
			}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
		})
	})
	Describe("DynamicConfiguration with Server2Cluster functional test", func() {
		It("SyncDeletion & SyncPolicyAlways", func() {
			rand.Seed(time.Now().Unix())
			randInt := rand.Int()
			dataId := fmt.Sprintf("randon-data-id-%d", randInt)
			content := fmt.Sprintf("content-%d", randInt)
			group := fmt.Sprintf("suite-test-group-server2cluster-%d", randInt)
			dc := v12.DynamicConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dc-suite-test-server2cluster",
					Namespace: dcTestNamespaceStr,
				},
				Spec: v12.DynamicConfigurationSpec{
					DataIds: []string{dataId},
					Strategy: v12.SyncStrategy{
						SyncPolicy:    v12.Always,
						SyncDeletion:  true,
						SyncDirection: v12.Server2Cluster,
					},
					NacosServer: v12.NacosServerConfiguration{
						Endpoint:  pointer.String(os.Getenv(nacosServerEndpoint)),
						Namespace: os.Getenv(nacosServerNamespace),
						Group:     group,
						AuthRef: &v1.ObjectReference{
							Name:       dcTestNacosCredentialName,
							APIVersion: "v1",
							Kind:       "Secret",
						},
					},
				},
			}
			By("create dataId and content in nacos server")
			createOrUpdateContentInNaocs(&dc, dataId, content)
			By("create a DynamicConfiguration with server2cluster sync direction")
			ensureDynamicConfigurationExist(&dc)
			By(fmt.Sprintf("check configmap: %s for content: %s", dc.Name, content))
			gomega.Eventually(checkConfigMapWithDataIdAndContent).
				WithArguments(dc.Name, dc.Namespace, dataId, content).
				WithTimeout(30 * time.Second).
				WithPolling(5 * time.Second).
				Should(gomega.BeTrue())
			newRandInt := rand.Int()
			newContent := fmt.Sprintf("content-%d", newRandInt)
			By("update dataId and content in nacos server")
			createOrUpdateContentInNaocs(&dc, dataId, newContent)
			By(fmt.Sprintf("check configmap: %s for changed content: %s", dc.Name, newContent))
			gomega.Eventually(checkConfigMapWithDataIdAndContent).
				WithArguments(dc.Name, dc.Namespace, dataId, newContent).
				WithTimeout(30 * time.Second).
				WithPolling(5 * time.Second).
				Should(gomega.BeTrue())
			By("delete dataId in nacos server")
			deleteDataIdInNaocs(&dc, dataId)
			By("delete DynamicConfiguration")
			gomega.Expect(k8sClient.Delete(ctx, &dc)).Should(gomega.Succeed())
			gomega.Eventually(func() bool {
				v := v12.DynamicConfiguration{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: dc.Name, Namespace: dc.Namespace}, &v)
				if err == nil {
					return false
				}
				gomega.Expect(errors.IsNotFound(err)).Should(gomega.BeTrue())
				return true
			}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
		})
		It("With spec.ObjectRef", func() {
			rand.Seed(time.Now().Unix() + 2)
			randInt := rand.Int()
			dataId := fmt.Sprintf("randon-data-id-%d", randInt)
			content := fmt.Sprintf("content-%d", randInt)
			group := fmt.Sprintf("suite-test-group-s2c-%d", randInt)
			cm := v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dc-suite-test-s2c-with-obj-ref",
					Namespace: dcTestNamespaceStr,
				},
			}
			dcWithObjRef := v12.DynamicConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dc-suite-test-s2c-with-obj-ref",
					Namespace: dcTestNamespaceStr,
				},
				Spec: v12.DynamicConfigurationSpec{
					DataIds: []string{dataId},
					Strategy: v12.SyncStrategy{
						SyncPolicy:    v12.Always,
						SyncDeletion:  true,
						SyncDirection: v12.Server2Cluster,
					},
					NacosServer: v12.NacosServerConfiguration{
						Endpoint:  pointer.String(os.Getenv(nacosServerEndpoint)),
						Namespace: os.Getenv(nacosServerNamespace),
						Group:     group,
						AuthRef: &v1.ObjectReference{
							Name:       dcTestNacosCredentialName,
							APIVersion: "v1",
							Kind:       "Secret",
						},
					},
					ObjectRef: &v1.ObjectReference{
						Name:       cm.Name,
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
				},
			}
			By("create dataId and content in nacos server")
			createOrUpdateContentInNaocs(&dcWithObjRef, dataId, content)
			By("create a DynamicConfiguration with server2cluster sync direction")
			ensureDynamicConfigurationExist(&dcWithObjRef)
			ensureConfigMapExist(&cm)
			By(fmt.Sprintf("check configmap: %s for content: %s", cm.Name, content))
			gomega.Eventually(checkConfigMapWithDataIdAndContent).
				WithArguments(cm.Name, dcWithObjRef.Namespace, dataId, content).
				WithTimeout(30 * time.Second).
				WithPolling(5 * time.Second).
				Should(gomega.BeTrue())
			By("delete dataId in nacos server")
			deleteDataIdInNaocs(&dcWithObjRef, dataId)
			By("delete DynamicConfiguration & ConfigMap")
			gomega.Expect(k8sClient.Delete(ctx, &dcWithObjRef)).Should(gomega.Succeed())
			gomega.Expect(k8sClient.Delete(ctx, &cm)).Should(gomega.Succeed())
			gomega.Eventually(func() bool {
				v := v12.DynamicConfiguration{}
				if err := k8sClient.Get(ctx, types.NamespacedName{Name: dcWithObjRef.Name, Namespace: dcWithObjRef.Namespace}, &v); err == nil {
					return false
				} else {
					gomega.Expect(errors.IsNotFound(err)).Should(gomega.BeTrue())
				}
				m := v1.ConfigMap{}
				if err := k8sClient.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, &m); err == nil {
					return false
				} else {
					gomega.Expect(errors.IsNotFound(err)).Should(gomega.BeTrue())
				}
				return true
			}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
		})
	})
})

func ensureNamespaceExist(namespace string) {
	testNamespace := v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	nn := types.NamespacedName{Name: namespace}
	err := k8sClient.Get(ctx, nn, &testNamespace)
	if err != nil {
		if errors.IsNotFound(err) {
			gomega.Expect(k8sClient.Create(ctx, &testNamespace)).Should(gomega.Succeed())
			gomega.Eventually(func() bool {
				return k8sClient.Get(ctx, nn, &testNamespace) == nil
			}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
		} else {
			gomega.Expect(err).ToNot(gomega.HaveOccurred(), fmt.Sprintf("get namespace %s failed", namespace))
		}
	}
	gomega.Expect(k8sClient.Update(ctx, &testNamespace)).Should(gomega.Succeed())
}

func ensureConfigMapExist(cm *v1.ConfigMap) {
	ensureObjectExist(cm, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	})
}

func ensureSecretExist(s *v1.Secret) {
	ensureObjectExist(s, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})
}

func ensureDynamicConfigurationExist(dc *v12.DynamicConfiguration) {
	ensureObjectExist(dc, schema.GroupVersionKind{
		Group:   "nacos.io",
		Version: "v1",
		Kind:    "DynamicConfiguration",
	})
}

func ensureObjectExist(obj client.Object, gvk schema.GroupVersionKind) {
	nn := types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}
	u := unstructured.Unstructured{}
	u.SetGroupVersionKind(gvk)
	err := k8sClient.Get(ctx, nn, &u)
	if err != nil {
		if !errors.IsNotFound(err) {
			gomega.Expect(err).ToNot(gomega.HaveOccurred(), fmt.Sprintf("get obj %s/%s failed", obj.GetNamespace(), obj.GetName()))
		}
		gomega.Expect(k8sClient.Create(ctx, obj)).Should(gomega.Succeed())
		gomega.Eventually(func() bool {
			return k8sClient.Get(ctx, nn, obj) == nil
		}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
	}
	gomega.Eventually(func() bool {
		err = k8sClient.Get(ctx, nn, &u)
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), fmt.Sprintf("get obj %s/%s failed", obj.GetNamespace(), obj.GetName()))
		obj.SetResourceVersion(u.GetResourceVersion())
		err = k8sClient.Update(ctx, obj)
		if err == nil {
			return true
		}
		gomega.Expect(errors.IsConflict(err)).Should(gomega.BeTrue(), fmt.Sprintf("update obj %s/%s failed", obj.GetNamespace(), obj.GetName()))
		return false
	}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
}

func cleanDynamicConfigurationTestResource() {
	secretNN := types.NamespacedName{Namespace: dcTestNamespaceStr, Name: dcTestNacosCredentialName}
	s := &v1.Secret{}
	err := k8sClient.Get(ctx, secretNN, s)
	if err != nil {
		if !errors.IsNotFound(err) {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get secret %s", secretNN.String()))
		}
	} else {
		gomega.Expect(k8sClient.Delete(ctx, s)).Should(gomega.Succeed())
	}

	nn := types.NamespacedName{Name: dcTestNamespaceStr}
	namespace := &v1.Namespace{}
	err = k8sClient.Get(ctx, nn, namespace)
	if err != nil {
		if !errors.IsNotFound(err) {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get namespace %s", nn.Namespace))
		}
	} else {
		gomega.Expect(k8sClient.Delete(ctx, namespace)).Should(gomega.Succeed())
	}
}

func checkDynamicConfigurationStatus(dc *v12.DynamicConfiguration) {
	gomega.Expect(dc != nil).Should(gomega.BeTrue())
	gomega.Eventually(func() bool {
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name:      dc.Name,
			Namespace: dc.Namespace,
		}, dc)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(dc.Status.Phase).NotTo(gomega.BeEquivalentTo("failed"), "dc status failed")
		return dc.Status.Phase == "succeed" && dc.Generation == dc.Status.ObservedGeneration
	}, time.Second*30, time.Second*5).Should(gomega.BeTrue())
}

func getContentByDataId(dc *v12.DynamicConfiguration) (string, bool) {
	configClient, err := impl.GetNacosClientBuilder().Build(&auth.DefaultNacosAuthProvider{Client: k8sClient}, dc)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	content, err := configClient.GetConfig(vo.ConfigParam{
		Group:  dc.Spec.NacosServer.Group,
		DataId: dc.Spec.DataIds[0],
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return content, len(content) > 0
}

func createOrUpdateContentInNaocs(dc *v12.DynamicConfiguration, dataId, content string) {
	configClient, err := impl.GetNacosClientBuilder().Build(&auth.DefaultNacosAuthProvider{Client: k8sClient}, dc)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	serverContent, err := configClient.GetConfig(vo.ConfigParam{
		Group:  dc.Spec.NacosServer.Group,
		DataId: dataId,
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	if len(serverContent) == 0 || serverContent != content {
		_, err := configClient.PublishConfig(vo.ConfigParam{
			Group:   dc.Spec.NacosServer.Group,
			DataId:  dataId,
			Content: content,
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
}

func deleteDataIdInNaocs(dc *v12.DynamicConfiguration, dataId string) {
	configClient, err := impl.GetNacosClientBuilder().Build(&auth.DefaultNacosAuthProvider{Client: k8sClient}, dc)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	_, err = configClient.DeleteConfig(vo.ConfigParam{
		Group:  dc.Spec.NacosServer.Group,
		DataId: dataId,
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func checkConfigMapWithDataIdAndContent(name, namespace, dataId, content string) bool {
	cm := v1.ConfigMap{}
	err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &cm)
	if err != nil {
		return false
	}
	c, ok := cm.Data[dataId]
	if !ok {
		return false
	}
	return c == content
}

package framework

import (
	"os"
	"testing"

	"github.com/aws/eks-anywhere/internal/pkg/api"
	anywherev1 "github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/executables"
	"github.com/aws/eks-anywhere/pkg/networkutils"
)

const (
	vsphereDatacenterVar        = "T_VSPHERE_DATACENTER"
	vsphereDatastoreVar         = "T_VSPHERE_DATASTORE"
	vsphereFolderVar            = "T_VSPHERE_FOLDER"
	vsphereNetworkVar           = "T_VSPHERE_NETWORK"
	vspherePrivateNetworkVar    = "T_VSPHERE_PRIVATE_NETWORK"
	vsphereResourcePoolVar      = "T_VSPHERE_RESOURCE_POOL"
	vsphereServerVar            = "T_VSPHERE_SERVER"
	vsphereSshAuthorizedKeyVar  = "T_VSPHERE_SSH_AUTHORIZED_KEY"
	vsphereStoragePolicyNameVar = "T_VSPHERE_STORAGE_POLICY_NAME"
	vsphereTemplateUbuntu118Var = "T_VSPHERE_TEMPLATE_UBUNTU_1_18"
	vsphereTemplateUbuntu119Var = "T_VSPHERE_TEMPLATE_UBUNTU_1_19"
	vsphereTemplateUbuntu120Var = "T_VSPHERE_TEMPLATE_UBUNTU_1_20"
	vsphereTemplateUbuntu121Var = "T_VSPHERE_TEMPLATE_UBUNTU_1_21"
	vsphereTemplateUbuntu122Var = "T_VSPHERE_TEMPLATE_UBUNTU_1_22"
	vsphereTemplateBR120Var     = "T_VSPHERE_TEMPLATE_BR_1_20"
	vsphereTemplateBR121Var     = "T_VSPHERE_TEMPLATE_BR_1_21"
	vsphereTlsInsecureVar       = "T_VSPHERE_TLS_INSECURE"
	vsphereTlsThumbprintVar     = "T_VSPHERE_TLS_THUMBPRINT"
	vsphereUsernameVar          = "EKSA_VSPHERE_USERNAME"
	vspherePasswordVar          = "EKSA_VSPHERE_PASSWORD"
	cidrVar                     = "T_VSPHERE_CIDR"
	privateNetworkCidrVar       = "T_VSPHERE_PRIVATE_NETWORK_CIDR"
	govcUrlVar                  = "GOVC_URL"
)

var requiredEnvVars = []string{
	vsphereDatacenterVar,
	vsphereDatastoreVar,
	vsphereFolderVar,
	vsphereNetworkVar,
	vspherePrivateNetworkVar,
	vsphereResourcePoolVar,
	vsphereServerVar,
	vsphereSshAuthorizedKeyVar,
	vsphereTemplateUbuntu118Var,
	vsphereTemplateUbuntu119Var,
	vsphereTemplateUbuntu120Var,
	vsphereTemplateUbuntu121Var,
	vsphereTemplateUbuntu122Var,
	vsphereTemplateBR120Var,
	vsphereTemplateBR121Var,
	vsphereTlsInsecureVar,
	vsphereTlsThumbprintVar,
	vsphereUsernameVar,
	vspherePasswordVar,
	cidrVar,
	privateNetworkCidrVar,
	govcUrlVar,
}

type VSphere struct {
	t              *testing.T
	fillers        []api.VSphereFiller
	clusterFillers []api.ClusterFiller
	cidr           string
	GovcClient     *executables.Govc
}

type VSphereOpt func(*VSphere)

func UpdateUbuntuTemplate118Var() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateUbuntu118Var, api.WithTemplateForAllMachines)
}

func UpdateUbuntuTemplate119Var() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateUbuntu119Var, api.WithTemplateForAllMachines)
}

func UpdateUbuntuTemplate120Var() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateUbuntu120Var, api.WithTemplateForAllMachines)
}

func UpdateUbuntuTemplate121Var() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateUbuntu121Var, api.WithTemplateForAllMachines)
}

func UpdateUbuntuTemplate122Var() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateUbuntu122Var, api.WithTemplateForAllMachines)
}

func UpdateBottlerocketTemplate121() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateBR121Var, api.WithTemplateForAllMachines)
}

func UpdateBottlerocketTemplate120() api.VSphereFiller {
	return api.WithStringFromEnvVar(vsphereTemplateBR120Var, api.WithTemplateForAllMachines)
}

func NewVSphere(t *testing.T, opts ...VSphereOpt) *VSphere {
	checkRequiredEnvVars(t, requiredEnvVars)
	c := buildGovc(t)
	v := &VSphere{
		t:          t,
		GovcClient: c,
		fillers: []api.VSphereFiller{
			api.WithStringFromEnvVar(vsphereDatacenterVar, api.WithDatacenter),
			api.WithStringFromEnvVar(vsphereDatastoreVar, api.WithDatastoreForAllMachines),
			api.WithStringFromEnvVar(vsphereFolderVar, api.WithFolderForAllMachines),
			api.WithStringFromEnvVar(vsphereNetworkVar, api.WithNetwork),
			api.WithStringFromEnvVar(vsphereResourcePoolVar, api.WithResourcePoolForAllMachines),
			api.WithStringFromEnvVar(vsphereServerVar, api.WithServer),
			api.WithStringFromEnvVar(vsphereSshAuthorizedKeyVar, api.WithSSHAuthorizedKeyForAllMachines),
			api.WithStringFromEnvVar(vsphereStoragePolicyNameVar, api.WithStoragePolicyNameForAllMachines),
			api.WithBoolFromEnvVar(vsphereTlsInsecureVar, api.WithTLSInsecure),
			api.WithStringFromEnvVar(vsphereTlsThumbprintVar, api.WithTLSThumbprint),
		},
	}

	v.cidr = os.Getenv(cidrVar)

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *VSphere) generateUniqueIp() string {
	ipgen := networkutils.NewIPGenerator(&networkutils.DefaultNetClient{})
	ip, err := ipgen.GenerateUniqueIP(v.cidr)
	if err != nil {
		v.t.Fatalf("Error getting unique IP for vsphere: %v", err)
	}
	return ip
}

func WithUbuntu122() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateUbuntu122Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Ubuntu),
		)
	}
}

func WithUbuntu121() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateUbuntu121Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Ubuntu),
		)
	}
}

func WithUbuntu120() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateUbuntu120Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Ubuntu),
		)
	}
}

func WithUbuntu119() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateUbuntu119Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Ubuntu),
		)
	}
}

func WithUbuntu118() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateUbuntu118Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Ubuntu),
		)
	}
}

func WithBottleRocket120() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateBR120Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Bottlerocket),
		)
	}
}

func WithBottleRocket121() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vsphereTemplateBR121Var, api.WithTemplateForAllMachines),
			api.WithOsFamilyForAllMachines(anywherev1.Bottlerocket),
		)
	}
}

func WithPrivateNetwork() VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers,
			api.WithStringFromEnvVar(vspherePrivateNetworkVar, api.WithNetwork),
		)
		v.cidr = os.Getenv(privateNetworkCidrVar)
	}
}

func WithVSphereWorkerNodeGroup(name string, workerNodeGroup *WorkerNodeGroup, fillers ...api.VSphereMachineConfigFiller) VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers, vSphereMachineConfig(name, fillers...))

		v.clusterFillers = append(v.clusterFillers, buildVSphereWorkerNodeGroupClusterFiller(name, workerNodeGroup))
	}
}

func WithVSphereFillers(fillers ...api.VSphereFiller) VSphereOpt {
	return func(v *VSphere) {
		v.fillers = append(v.fillers, fillers...)
	}
}

func (v *VSphere) Name() string {
	return "vsphere"
}

func (v *VSphere) Setup() {}

func (v *VSphere) CustomizeProviderConfig(file string) []byte {
	return v.customizeProviderConfig(file, v.fillers...)
}

func (v *VSphere) customizeProviderConfig(file string, fillers ...api.VSphereFiller) []byte {
	providerOutput, err := api.AutoFillVSphereProvider(file, fillers...)
	if err != nil {
		v.t.Fatalf("Error customizing provider config from file: %v", err)
	}
	return providerOutput
}

func (v *VSphere) WithProviderUpgrade(fillers ...api.VSphereFiller) ClusterE2ETestOpt {
	return func(e *ClusterE2ETest) {
		e.ProviderConfigB = v.customizeProviderConfig(e.ClusterConfigLocation, fillers...)
	}
}

func (v *VSphere) WithProviderUpgradeGit(fillers ...api.VSphereFiller) ClusterE2ETestOpt {
	return func(e *ClusterE2ETest) {
		e.ProviderConfigB = v.customizeProviderConfig(e.clusterConfigGitPath(), fillers...)
	}
}

func (v *VSphere) WithNewVSphereWorkerNodeGroup(name string, workerNodeGroup *WorkerNodeGroup, fillers ...api.VSphereMachineConfigFiller) ClusterE2ETestOpt {
	return func(e *ClusterE2ETest) {
		e.ProviderConfigB = v.customizeProviderConfig(e.ClusterConfigLocation, vSphereMachineConfig(name, fillers...))
		var err error
		// Using the ClusterConfigB instead of file in disk since it might have already been updated but not written to disk
		e.ClusterConfigB, err = api.AutoFillClusterFromYaml(e.ClusterConfigB, buildVSphereWorkerNodeGroupClusterFiller(name, workerNodeGroup))
		if err != nil {
			e.T.Fatalf("Error filling cluster config: %v", err)
		}
	}
}

func (v *VSphere) ClusterConfigFillers() []api.ClusterFiller {
	value, ok := os.LookupEnv(ClusterIPPoolEnvVar)
	var clusterIP string
	var err error
	if ok && value != "" {
		clusterIP, err = PopIPFromEnv(ClusterIPPoolEnvVar)
		if err != nil {
			v.t.Fatalf("failed to pop cluster ip from test environment: %v", err)
		}
	} else {
		clusterIP = v.generateUniqueIp()
	}

	v.clusterFillers = append(v.clusterFillers, api.WithControlPlaneEndpointIP(clusterIP))
	return v.clusterFillers
}

func RequiredVsphereEnvVars() []string {
	return requiredEnvVars
}

func vSphereMachineConfig(name string, fillers ...api.VSphereMachineConfigFiller) api.VSphereFiller {
	f := make([]api.VSphereMachineConfigFiller, 0, len(fillers)+6)
	// Need to add these because at this point the default fillers that assign these
	// values to all machines have already ran
	f = append(f,
		api.WithVSphereMachineDefaultValues(),
		api.WithDatastore(os.Getenv(vsphereDatastoreVar)),
		api.WithFolder(os.Getenv(vsphereFolderVar)),
		api.WithResourcePool(os.Getenv(vsphereResourcePoolVar)),
		api.WithStoragePolicyName(os.Getenv(vsphereStoragePolicyNameVar)),
		api.WithSSHKey(os.Getenv(vsphereSshAuthorizedKeyVar)),
	)
	f = append(f, fillers...)

	return api.WithMachineConfig(name, f...)
}

func buildVSphereWorkerNodeGroupClusterFiller(machineConfigName string, workerNodeGroup *WorkerNodeGroup) api.ClusterFiller {
	// Set worker node group ref to vsphere machine config
	workerNodeGroup.MachineConfigKind = anywherev1.VSphereMachineConfigKind
	workerNodeGroup.MachineConfigName = machineConfigName
	return workerNodeGroup.clusterFiller()
}

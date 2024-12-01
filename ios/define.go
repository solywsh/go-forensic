package ios

const (
	ApplicationPath       = "/private/var/mobile/Containers/Data/Application"
	AppGroupPath          = "/private/var/mobile/Containers/Shared/AppGroup"
	ContainerMetaDataFile = ".com.apple.mobile_container_manager.metadata.plist"

	ActiveFileName = "active.tar"
)

var focusPathList = []string{
	ApplicationPath,
	AppGroupPath,
}

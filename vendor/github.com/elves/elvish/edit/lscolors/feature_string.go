// Code generated by "stringer -type=feature -output=feature_string.go"; DO NOT EDIT.

package lscolors

import "strconv"

const _feature_name = "featureInvalidfeatureOrphanedSymlinkfeatureSymlinkfeatureMultiHardLinkfeatureNamedPipefeatureSocketfeatureDoorfeatureBlockDevicefeatureCharDevicefeatureWorldWritableStickyDirectoryfeatureWorldWritableDirectoryfeatureStickyDirectoryfeatureDirectoryfeatureCapabilityfeatureSetuidfeatureSetgidfeatureExecutablefeatureRegular"

var _feature_index = [...]uint16{0, 14, 36, 50, 70, 86, 99, 110, 128, 145, 180, 209, 231, 247, 264, 277, 290, 307, 321}

func (i feature) String() string {
	if i < 0 || i >= feature(len(_feature_index)-1) {
		return "feature(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _feature_name[_feature_index[i]:_feature_index[i+1]]
}

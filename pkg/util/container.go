/*
 * Copyright 2022 The Koordinator Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/koordinator-sh/koordinator/apis/extension"
	"github.com/koordinator-sh/koordinator/pkg/util/system"
)

// @parentDir kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/
// @return kubepods.slice/kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/****.scope
func GetContainerCgroupPathWithKube(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerDir, err := system.CgroupPathFormatter.ContainerDirFn(c)
	if err != nil {
		return "", err
	}
	return path.Join(
		GetPodCgroupDirWithKube(podParentDir),
		containerDir,
	), nil
}

func GetContainerCgroupCPUAcctProcStatPath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.CpuacctStat), nil
}

func GetContainerCgroupMemStatPath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.MemStat), nil
}

func GetContainerCgroupCPUStatPath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.CPUStat), nil
}

func GetContainerCgroupMemLimitPath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.MemoryLimit), nil
}

func GetContainerBEMilliCPURequest(c *corev1.Container) int64 {
	if cpuRequest, ok := c.Resources.Requests[extension.BatchCPU]; ok {
		return cpuRequest.Value()
	}
	return -1
}

func GetContainerBEMilliCPULimit(c *corev1.Container) int64 {
	if cpuRequest, ok := c.Resources.Limits[extension.BatchCPU]; ok {
		return cpuRequest.Value()
	}
	return -1
}

func GetContainerBEMemoryByteLimit(c *corev1.Container) int64 {
	if memLimit, ok := c.Resources.Limits[extension.BatchMemory]; ok {
		return memLimit.Value()
	}
	return -1
}

func GetContainerCgroupCPUSharePath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.CPUShares), nil
}

func GetContainerCgroupCFSPeriodPath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.CPUCFSPeriod), nil
}

func GetContainerCgroupCFSQuotaPath(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	containerPath, err := GetContainerCgroupPathWithKube(podParentDir, c)
	if err != nil {
		return "", err
	}
	return system.GetCgroupFilePath(containerPath, system.CPUCFSQuota), nil
}

func GetContainerCurCPUShare(podParentDir string, c *corev1.ContainerStatus) (int64, error) {
	cgroupPath, err := GetContainerCgroupCPUSharePath(podParentDir, c)
	if err != nil {
		return 0, err
	}
	rawContent, err := ioutil.ReadFile(cgroupPath)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(rawContent)), 10, 64)
}

func GetContainerCurCFSPeriod(podParentDir string, c *corev1.ContainerStatus) (int64, error) {
	cgroupPath, err := GetContainerCgroupCFSPeriodPath(podParentDir, c)
	if err != nil {
		return 0, err
	}
	rawContent, err := ioutil.ReadFile(cgroupPath)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(rawContent)), 10, 64)
}

func GetContainerCurCFSQuota(podParentDir string, c *corev1.ContainerStatus) (int64, error) {
	cgroupPath, err := GetContainerCgroupCFSQuotaPath(podParentDir, c)
	if err != nil {
		return 0, err
	}
	rawContent, err := ioutil.ReadFile(cgroupPath)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(rawContent)), 10, 64)
}

func GetContainerCurMemLimitBytes(podParentDir string, c *corev1.ContainerStatus) (int64, error) {
	cgroupPath, err := GetContainerCgroupMemLimitPath(podParentDir, c)
	if err != nil {
		return 0, err
	}
	rawContent, err := ioutil.ReadFile(cgroupPath)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(rawContent)), 10, 64)
}

func FindContainerIdAndStatusByName(status *corev1.PodStatus, name string) (string, *corev1.ContainerStatus, error) {
	allStatuses := status.InitContainerStatuses
	allStatuses = append(allStatuses, status.ContainerStatuses...)
	for _, container := range allStatuses {
		if container.Name == name && container.ContainerID != "" {
			_, cID, err := ParseContainerId(container.ContainerID)
			if err != nil {
				return "", nil, err
			}
			return cID, &container, nil
		}
	}
	return "", nil, fmt.Errorf("unable to find ID for container with name %v in pod status (it may not be running)", name)
}

func ParseContainerId(data string) (cType, cID string, err error) {
	// Trim the quotes and split the type and ID.
	parts := strings.Split(strings.Trim(data, "\""), "://")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid container ID: %q", data)
		return
	}
	cType, cID = parts[0], parts[1]
	return
}
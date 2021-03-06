// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"strconv"

	"github.com/spf13/cobra"
)

var deployFlagVerifyingChain uint8

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a smart contract (Testimonium or Ethash) on the specified blockchain",
	Long: `Deploys a smart contract (Testimonium or Ethash) on the specified blockchain`,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deployCmd.PersistentFlags().Uint8VarP(&deployFlagVerifyingChain, "verifying", "v", 1, "The blockchain to which the smart contract is deployed")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func updateChainsConfig(deployedAddress common.Address, chainId uint8, key string) {
	chainsConfig := viper.Get("chains").(map[string]interface {})
	deployChainConfig := chainsConfig[strconv.FormatUint(uint64(chainId), 10)].(map[string]interface{})
	deployChainConfig[key] = deployedAddress.String()

	chainsConfig[strconv.FormatUint(uint64(chainId), 10)] = deployChainConfig
	viper.Set("chains", chainsConfig)

	_ = viper.WriteConfig()
}

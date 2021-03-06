// This file contains logic executed if the command "verify transaction" is typed in.
// Authors: Marten Sigwart, Philipp Frauenthaler

package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pantos-io/go-ethrelay/testimonium"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var noOfConfirmations uint8

// verifyTransactionCmd represents the transaction command
var verifyTransactionCmd = &cobra.Command{
	Use:   "transaction [txHash]",
	Short: "Verifies a transaction",
	Long: `Verifies a transaction from the target chain on the verifying chain

Behind the scene, the command queries the transaction with the specified hash ('txHash') from the target chain.
It then generates a Merkle Proof contesting the existence of the transaction within a specific block.
This information gets sent to the verifying chain, where not only the existence of the block but also the Merkle Proof are verified`,
	Aliases: []string{"tx"},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txHash := common.HexToHash(args[0])

		testimoniumClient = createTestimoniumClient()

		rlpHeader, rlpEncodedTx, path, rlpEncodedProofNodes, err := testimoniumClient.GenerateMerkleProofForTx(txHash, verifyFlagSrcChain)
		if err != nil {
			log.Fatal("Failed to generate Merkle Proof: " + err.Error())
		}

		// TODO: this produces a merkle proof for the transaction and does not verify the transaction
		//  maybe it is better to introduce a new command for this behaviour as it is quite confusing to
		//  call verifyTransaction and no transaction is verified
		if jsonFlag {
			hexEncodedTxHash := make([]byte, hex.EncodedLen(len(txHash)))
			hex.Encode(hexEncodedTxHash, txHash[:])

			writeMerkleProofAsJson(hexEncodedTxHash, rlpHeader, rlpEncodedTx, path, rlpEncodedProofNodes)

			fmt.Printf("Wrote merkle proof to 0x%s.json\n", hexEncodedTxHash)

			return
		}

		feesInWei, err := testimoniumClient.GetRequiredVerificationFee(verifyFlagDestChain)
		if err != nil {
			log.Fatal(err)
		}

		testimoniumClient.VerifyMerkleProof(feesInWei, rlpHeader, testimonium.VALUE_TYPE_TRANSACTION, rlpEncodedTx, path,
			rlpEncodedProofNodes, noOfConfirmations, verifyFlagDestChain)
	},
}

func writeMerkleProofAsJson(fileName []byte, rlpHeader []byte, tx []byte, path []byte, nodes []byte) {
	f, err := os.Create(fmt.Sprintf("./0x%s.json", fileName))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	hexEncodedRlpHeader := make([]byte, hex.EncodedLen(len(rlpHeader)))
	hex.Encode(hexEncodedRlpHeader, rlpHeader)

	hexEncodedTx := make([]byte, hex.EncodedLen(len(tx)))
	hex.Encode(hexEncodedTx, tx)

	hexEncodedPath := make([]byte, hex.EncodedLen(len(path)))
	hex.Encode(hexEncodedPath, path)

	hexEncodedNodes := make([]byte, hex.EncodedLen(len(nodes)))
	hex.Encode(hexEncodedNodes, nodes)

	_, err = fmt.Fprint(f, "{\n")
	_, err = fmt.Fprintf(f, "  \"rlpHeader\": \"0x%s\",\n", hexEncodedRlpHeader)
	_, err = fmt.Fprintf(f, "  \"rlpEncodedTx\": \"0x%s\",\n", hexEncodedTx)
	_, err = fmt.Fprintf(f, "  \"path\": \"0x%s\",\n", hexEncodedPath)
	_, err = fmt.Fprintf(f, "  \"rlpEncodedNodes\": \"0x%s\"\n", hexEncodedNodes)
	_, err = fmt.Fprint(f, "\n}")

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	verifyCmd.AddCommand(verifyTransactionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyTransactionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	verifyTransactionCmd.Flags().Uint8VarP(&noOfConfirmations, "confirmations", "c", 4, "Number of block confirmations")
	verifyTransactionCmd.Flags().BoolVar(&jsonFlag, "json", false, "save merkle proof to a json file")
}

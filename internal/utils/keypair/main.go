package keypair

import (
	"github.com/stellar/go/keypair"
	"log"
)

type Keypair struct {
	public string
	private string
}

func Generate() *Keypair {
	pair, err := keypair.Random()

	if err != nil {
		log.Fatal(err)
	}

	kp := &Keypair{
		private: pair.Seed(),
		public: pair.Address(),
	}

	log.Println("Generated new keypair:", kp.public, kp.private)
	// SAV76USXIJOBMEQXPANUOQM6F5LIOTLPDIDVRJBFFE2MDJXG24TAPUU7
	// GCFXHS4GXL6BVUCXBWXGTITROWLVYXQKQLF4YH5O5JT3YZXCYPAFBJZB
	return kp
}
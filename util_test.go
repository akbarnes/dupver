package dupver

import "testing"

func TestReadFilters(test *testing.T) {
   filters, err := ReadFilters() 

    if err != nil {
        test.Errorf("Error reading filters from file")
    }

    if len(filters) != 2 {
        test.Errorf("Filter length of %d != 2", len(filters))
    }

   if filters[0] != "filter1" {
        test.Errorf("Filter 1 value of %s != filter1", filters[0])
    }

   if filters[1] != "filter2" {
        test.Errorf("Filter 2 value of %s != filter1", filters[1])
    }

}

func TestRandHexString(test *testing.T) {
    n := 32
    s := RandHexString(n)

    if len(s) != n {
        test.Errorf("Random string length %d != %d", len(s), n)
    }

    hexCharSet := map[rune]bool{}

    for _, c := range HexChars {
        hexCharSet[c] = true
    }

    for _, c := range s {
        if _, ok := hexCharSet[c]; !ok {
            test.Errorf("Character %c isn't a valid hex character", c)
        }
    }
}

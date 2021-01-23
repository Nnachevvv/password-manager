package cmd_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/nnachevv/passmag/cmd"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = Describe("Change", func() {
	var (
		c      *expect.Console
		state  *vt10x.State
		err    error
		path   string
		getCmd *cobra.Command
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		getCmd = cmd.NewGetCmd(terminal.Stdio{c.Tty(), c.Tty(), c.Tty()})

		getCmd.SetArgs([]string{})
		getCmd.SetOut(&stdOut)
		getCmd.SetErr(&stdErr)

		path, err = tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())

		viper.Set("password.path", path)
		viper.Set("PASS_SESSION", "MRfbladUgDxLHvVWbxUjQUiZQykqiNcK")
	})

	Context("get existing password", func() {
		It("should return password", func() {
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for which you want to get your password:")
				c.SendLine("exist@mail.com")
				c.ExpectEOF()
			}()

			err = getCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(expect.StripTrailingEmptyLines(string(stdOut.String()))).To(Equal("gMdLasZIGAEmDSCprqFkZQSAnjzeZzUP"))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("get non-existing password", func() {
		It("should throw password", func() {
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for which you want to get your password:")
				c.SendLine("nonexist@mail.com")
				c.ExpectEOF()
			}()

			err = getCmd.Execute()
			Expect(err).To(Equal(errors.New("this name not exist in your password manager")))
			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("pass wrong master password", func() {
		It("throw failed to find this name error", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("wrong")
				c.ExpectEOF()

			}()
			err = getCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
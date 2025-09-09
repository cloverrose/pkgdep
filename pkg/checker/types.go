//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE --destination=mock_test.go
package checker

type recorder interface {
	RecordUsage(frm, to string)
}

package convert

// func Test_convertFile(t *testing.T) {
//   type args struct {
//     files []string
//   }
//   tests := []struct {
//     name    string
//     args    args
//     wantOut string
//     wantErr bool
//   }{
//     {
//       name: "convert file",
//       args: args{
//         files: []string{"../../../test/v1/data-source/data-source.yaml"},
//       },
//       wantOut: `foo`,
//       wantErr: false,
//     },
//   }
//   for _, tt := range tests {
//     t.Run(tt.name, func(t *testing.T) {
//       out := &bytes.Buffer{}
//       if err := convertFile(out, tt.args.files); (err != nil) != tt.wantErr {
//         t.Errorf("convertFile() error = %v, wantErr %v", err, tt.wantErr)
//         return
//       }
//       if gotOut := out.String(); gotOut != tt.wantOut {
//         t.Errorf("convertFile() got: \n%v, want %v", gotOut, tt.wantOut)
//       }
//     })
//   }
// }

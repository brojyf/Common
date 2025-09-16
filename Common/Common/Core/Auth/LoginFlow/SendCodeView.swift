//
//  SendCodeView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct SendCodeView: View {
    
    @EnvironmentObject var vm: AuthVM
    
    let scene: AuthScene
    @State private var email: String = ""
    
    var body: some View {
        VStack {
            InputField("email", text: $email)
            Button("Send Code"){
                vm.requestCodeWithRouter(email: email, scene: scene)
            }
        }
        .padding()
        .navigationTitle(Text(scene == .signup ? "Sign up" : "Reset Password"))
        .alert(isPresented: $vm.hasError){
            Alert(
                title: Text("Error"),
                message: Text(vm.errorMsg ?? "Unknown Error"),
                dismissButton: .default(Text("OK")){
                    vm.dismissError()
                }
            )
        }
    }
}

#Preview {
    let dev = dev.loggedOut()
    NavigationStack {
        SendCodeView(scene: .signup)
    }
    .environmentObject(dev.authVM)
}
#Preview {
    let dev = dev.loggedOut()
    NavigationStack {
        SendCodeView(scene: .resetPassword)
    }
    .environmentObject(dev.authVM)
}

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
                vm.requestCode(email: email, scene:"signup")
                vm.requestCodeWithRouter(email: email, scene: scene)            }
        }
        .padding()
        .navigationTitle(Text(scene == .signup ? "Sign up" : "Reset Password"))
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

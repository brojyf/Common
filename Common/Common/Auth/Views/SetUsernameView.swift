//
//  SetUsernameView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct SetUsernameView: View {
    
    @EnvironmentObject var authVM: AuthVM
    
    @State private var username: String = ""
    
    var body: some View {
        VStack {
            InputField("Username", text: $username)
            Button("Submit"){
                authVM.setUsernameFirstTime()
            }
        }
        .padding()
        .navigationBarTitle("Set Username")
    }
}

#Preview {
    let dev = dev.loggedOut()
    NavigationStack {
        SetUsernameView()
    }
    .environmentObject(dev.authVM)
}

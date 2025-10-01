//
//  AccountView.swift
//  Common
//
//  Created by 江逸帆 on 9/30/25.
//

import SwiftUI

struct AccountView: View {
    
    @EnvironmentObject var vm: AuthVM
    
    var body: some View {
        VStack {
            HStack {
                Text("This is Account View")
            }
            
            Button("Logout"){
                withAnimation(.spring){
                    vm.logout()
                }
            }
        }
    }
}

#Preview {
    let dev = dev.loggedIn()
    AccountView()
        .environmentObject(dev.authVM)
}

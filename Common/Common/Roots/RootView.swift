//
//  RootView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct RootView: View {
    
    @EnvironmentObject var session: SessionStore
    @EnvironmentObject var authVM: AuthVM
    
    var body: some View {
        Group {
            if session.isLoggedIn {
                NavigationStack{
                    MainAppRoot()
                }
            } else {
                LoginFlowRoot()
            }
        }
    }
}

#Preview {
    let dev = dev.loggedOut()
    NavigationStack {
        RootView()
    }
    .environmentObject(dev.session)
    .environmentObject(dev.authVM)
}

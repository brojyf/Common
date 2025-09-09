//
//  MainAppRoot.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct MainAppRoot: View {
    @EnvironmentObject var session: SessionStore
    
    var body: some View {
        Text("TabViews")
        Button("Logout"){
            withAnimation(.spring){
                session.logout()
            }
            
        }
    }
}

#Preview {
    let dev = dev.loggedIn()
    MainAppRoot()
        .environmentObject(dev.session)
}

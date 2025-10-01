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
        
        TabView{
            TodoView()
                .tabItem{
                    Label("Todos", systemImage: "list.bullet")
                }
            
            BillView()
                .tabItem{
                    Label("Bills", systemImage: "list.bullet")
                }
            
            AddView()
                .tabItem{
                    Label("Add", systemImage: "list.bullet")
                }
            
            FriendsView()
                .tabItem{
                    Label("Friends", systemImage: "list.bullet")
                }
            
            AccountView()
                .tabItem{
                    Label("Account", systemImage: "list.bullet")
                }
        }
    }
}


#Preview {
    let dev = dev.loggedIn()
    MainAppRoot()
        .environmentObject(dev.authVM)
}

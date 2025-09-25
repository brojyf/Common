//
//  KCManager.swift
//  Common
//
//  Created by 江逸帆 on 9/24/25.
//

import Foundation

enum KCKey: String {
    case codeID
    case ott
}

final class KCManager {
    
    @discardableResult
    static func save(_ key: KCKey, _ value: String) -> Bool {
        
        guard let data = value.data(using: .utf8) else { return false }
        let queryDelete: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key.rawValue
        ]
        SecItemDelete(queryDelete as CFDictionary)
                
        let queryAdd: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key.rawValue,
            kSecValueData as String: data
        ]
                
        let status = SecItemAdd(queryAdd as CFDictionary, nil)
        return status == errSecSuccess
    }
    
    @discardableResult
    static func load(_ key: KCKey) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key.rawValue,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]
            
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        
        if status == errSecSuccess,
            let data = result as? Data,
            let value = String(data: data, encoding: .utf8) {
                return value
        }
        return nil
    }
    
    @discardableResult
        static func delete(_ key: KCKey) -> Bool {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key.rawValue
        ]
        
        let status = SecItemDelete(query as CFDictionary)
        return status == errSecSuccess
    }
}

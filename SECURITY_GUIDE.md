# Security Guide for Home Automation System

## 🔐 Authorization Implementation Summary

Your home automation system now includes **Basic HTTP Authentication** with the following features:

### ✅ Implemented Security Features

1. **Basic Authentication**
   - Simple, reliable HTTP Basic Auth
   - Built-in browser support
   - Two user roles: `admin` and `user`

2. **Rate Limiting**
   - 5 attempts per 15 minutes per IP
   - Prevents brute force attacks

3. **Role-Based Access Control**
   - `admin`: Full system access
   - `user`: Control and monitoring access
   - `guest`: Read-only access (if needed)

4. **Security Headers**
   - HSTS (HTTP Strict Transport Security)
   - Content Security Policy (CSP)
   - X-Frame-Options (Clickjacking protection)
   - X-Content-Type-Options (MIME sniffing protection)

5. **Audit Logging**
   - All user actions are logged
   - Username and role tracking

## 🚀 Quick Start

### Default Credentials
```
Admin User:
  Username: admin
  Password: admin123!

Regular User:
  Username: family  
  Password: family123!
```

**⚠️ CRITICAL: Change these passwords immediately!**

### Running the System
```bash
go run *.go
```

The system will show:
```
=== SYSTEM AUTORYZACJI AKTYWNY ===
Serwer uruchomiony na porcie 8080
Dostęp wymaga autoryzacji Basic Auth
```

## 🔒 Production Security Checklist

### Before Going Public:

#### 1. **Change Default Passwords**
```go
// In basic_auth.go, modify NewBasicAuthManager():
auth.AddUser("admin", "YOUR_STRONG_PASSWORD", "admin")
auth.AddUser("family", "YOUR_FAMILY_PASSWORD", "user")
```

#### 2. **Enable HTTPS**
```bash
# Generate SSL certificate (Let's Encrypt recommended)
certbot certonly --standalone -d yourdomain.com

# Or for testing, generate self-signed:
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
```

#### 3. **Configure Reverse Proxy (Recommended)**

**Nginx Configuration:**
```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /path/to/server.crt;
    ssl_certificate_key /path/to/server.key;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

#### 4. **Firewall Configuration**
```bash
# Ubuntu/Debian
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP (for Let's Encrypt)
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# Block direct access to Go app port
sudo ufw deny 8080
```

#### 5. **Environment Variables**
```bash
# Create .env file
export HOME_AUTOMATION_ADMIN_PASS="your_secure_password"
export HOME_AUTOMATION_USER_PASS="your_user_password"
export TLS_CERT_PATH="/path/to/cert.pem"
export TLS_KEY_PATH="/path/to/key.pem"
```

## 🛡️ Security Best Practices

### Network Security
- **Use VPN**: Consider VPN access instead of direct internet exposure
- **Port Forwarding**: Only forward necessary ports (443 for HTTPS)
- **IP Whitelisting**: Restrict access to known IP ranges if possible

### Application Security
- **Regular Updates**: Keep Go and dependencies updated
- **Log Monitoring**: Monitor access logs for suspicious activity
- **Backup Strategy**: Regular backups of configuration and logs

### Advanced Options

#### 1. **Two-Factor Authentication (Future Enhancement)**
```go
// Add to BasicUser struct:
type BasicUser struct {
    Username     string
    PasswordHash string
    Role         string
    TOTPSecret   string  // For 2FA
    BackupCodes  []string
    LastAccess   time.Time
}
```

#### 2. **IP Whitelisting**
```go
func (ba *BasicAuthManager) IPWhitelistMiddleware(allowedIPs []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            clientIP := getClientIP(r)
            
            allowed := false
            for _, ip := range allowedIPs {
                if clientIP == ip {
                    allowed = true
                    break
                }
            }
            
            if !allowed {
                http.Error(w, "Access denied", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 📊 Monitoring and Logging

### Access Logs
The system logs all user actions:
```
Użytkownik admin (admin) zmienia moc na 60%
```

### Security Events to Monitor
- Multiple failed login attempts
- Access from new IP addresses
- Unusual usage patterns
- System configuration changes

## 🔧 Troubleshooting

### Common Issues

1. **Browser keeps asking for credentials**
   - Check username/password
   - Clear browser cache
   - Verify server is running

2. **"Too many requests" error**
   - Wait 15 minutes for rate limit reset
   - Check if IP is correct

3. **HTTPS certificate issues**
   - Verify certificate paths
   - Check certificate expiration
   - Ensure proper file permissions

## 📱 Mobile Access

Basic Auth works perfectly with mobile browsers and apps:
- **iOS Safari**: Built-in support
- **Android Chrome**: Built-in support
- **Home Assistant**: Compatible
- **Custom Apps**: Easy HTTP Basic Auth integration

## 🎯 Why Basic Auth is Perfect for Home Automation

1. **Simplicity**: No complex token management
2. **Reliability**: Works everywhere, always
3. **Security**: When used with HTTPS, very secure
4. **Compatibility**: Works with all clients and tools
5. **No JavaScript Required**: Works even if JS is disabled
6. **Stateless**: No server-side session management needed

## 🚨 Emergency Access

If you get locked out:
1. SSH to your server
2. Restart the Go application
3. Rate limits will reset
4. Or modify the code to temporarily disable auth

## 📞 Support

For security questions or issues:
1. Check logs: `journalctl -u your-service-name`
2. Test locally first: `curl -u admin:password http://localhost:8080/api/status`
3. Verify network connectivity and firewall rules

---

**Remember**: Security is a process, not a destination. Regularly review and update your security measures!

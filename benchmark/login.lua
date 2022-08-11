-- 登录系统，并将cookie保存到文件cookie.txt

wrk.method = "POST"  
wrk.body   = '{"username": "jinkesi", "password": "rootpwd"}'  
wrk.headers["Content-Type"] = "application/json"

index = 1
io.open("cookie.txt","w"):close()
file = io.open("cookie.txt", "a")	


function done() 
  print("done...")
  io.close(file)
end


request = function() 
    local userPrefix = "testu_"
    local method = "POST"
    local headers = { }
    body ='{"username": "' .. userPrefix .. index .. '", "password": "rootpwd"}'   
    headers['Content-Type'] = "application/json"
    index = index + 1
    return wrk.format(method, nil, headers, body)
  end

function getCookie(cookies, name)  
  local start = string.find(cookies, name .. "=")  
  
  if start == nil then  
    return nil  
  end  
  
  return string.sub(cookies, start + #name + 1, string.find(cookies, ";", start) - 1)  
end  

function response(status, headers, body)  
  local sessionId = getCookie(headers["Set-Cookie"], "sessionId")  
    
  if sessionId ~= nil then  
    print('sessionId: ' .. sessionId .. "\n")
    io.output(file) 
    io.write(sessionId .. '\n') 
    -- wrk.headers["Cookie"] = "token=" .. token  
  end 
end  

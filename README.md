<div>
  <h1>Read_and_Dispatch_site</h1>
  <p>I made site where you can publist your photos with description. For an API I used golang Gin. I made account registration, log in and authentication too</p>
</div>
<br>
<div>
  <h1>How to start using?</h1>
  <p>
    Be sure, that you have already installed <i>MySql DB</i> server. <br>
    It should work on the default port <strong>3306</strong> <br>
    (you can change port, but then you should change port inside <b>main.go</b>) <br>
    <br>
    After connecting to the DB enter <b>Backend/DB</b> directory and run <b>Users_DB.sql</b> script to create Users_DB.<br>
    This DB will store User's credentials and ids of dispatches in tables <strong>Users</strong> and <strong>Users_Dispatches</strong><br>
    (All Users have unique UserName and UserEmail. All dispatches have unique ID)<br>
    <br>
    For running backend you should install <i>Golang</i> <br>
    For running frontend you should use <i>Vite.js</i> with <i>React</i> (without <i>Typescript</i>)<br>
    <br>
    To run backend you should enter <b>Backend/Golang</b> directory with <b>main.go</b> and <b>go.mod</b> and then write in console: <br>
    <pre>
       $  go get . 
       $  go run . </pre>
    It will install all <i>Go</i> dependencies and run backend on <strong>http://localhost:8080</strong> <br>
    (you can change it inside <b>main.go</b>, but then you should change it inside frontend's <b>main.jsx</b>) <br>
    <br>
    To run frontend you should enter <b>Frontend</b> directory with <b>package.json</b> and other files, then write in console:
    <pre>
       $  npm install  
       $  npm run dev </pre>
    It will install all <i>React</i> dependencies and run frontend on <strong>http://localhost:5173</strong> <br>
    (you can change server's ip and port inside <b>vite.config.js</b>, but then you should change it inside backend's <b>main.go</b>) <br>
    <br>
    <strong><em>After all you will get working website where you can post and read posts!</em></strong>
    <br>
  </p>
  <br>
  <h3>Notice:</h3>
    <h4>
      <ul>
        <li>
          Password and other User's credentials are sended to backend and stored inside DB without any hashing or encryption. <br>
          It could cause troubles with safety. 
        </li>
        <br>
        <li>
          Function getDispatches() in <b>Frontend/Dispatcher/Dispatcher.jsx</b> get all ids of dispatches from DB. <br>
          It could cause troubles with optimization and page rendering if there would be a lot of dispatches. 
        </li>
      </ul>
    </h4>
</div>
<br>
<div align="center">
  <h1>Screenshots</h1>
  <img src="https://github.com/user-attachments/assets/91cdb310-ebae-4961-826b-f5b85036b7c3" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/6c7391a4-d4d5-4a9f-9824-c6a35a4eb839" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/18b76a99-28ec-4875-babc-91932f05b10d" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/22775c42-7d55-46bf-ae31-eee12c1290f6" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/f4045571-bad3-4919-8ad5-e80fd98b6901" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/e4e7a1a6-1ca9-426f-a83d-d5788a98d947" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/879a3889-4696-41a2-b01b-49edd5c7fb97" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/5dbef3d3-95c3-4da1-8d23-3db72a6de859" height="165vw" width="330vw">
  <img src="https://github.com/user-attachments/assets/a7885475-bcae-43e4-8ae6-8d8e6d365a06" height="165vw" width="330vw">
</div>

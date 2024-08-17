import { useEffect, useRef, useState } from 'react';
import './Dispatch.css'
import { useParams } from "react-router-dom";
import { useAuth } from '../Providers/AuthProvider';
import { API_ADRESS } from '../main';

export function Dispatch() { // Page with dispatch data (image, title, description and text)
  const { token, setToken } = useAuth();
  const [canEdit, setCanEdit] = useState(false);
  const { id } = useParams()

  const image = useRef(null)
  const title = useRef(null)
  const description = useRef(null)
  const text = useRef(null)

  async function createDispatch(token) { // Request to backend to create new dispatch
    if (title.current.value == "") { // All dispatches should have title
      alert("Write title!")
    } else {

      var dispatch = new FormData()

      dispatch.append('File', image.current.files[0])
      dispatch.append('Title', title.current.value)
      dispatch.append('Description', description.current.value)
      dispatch.append('Text', text.current.value)


      await fetch(`${API_ADRESS}/createDispatch`, {
        method: "POST",
        headers:
        {
          "Authorization": token,
        },
        body: dispatch
      }).then(response => {

        if (response.status == 201) {

          return response.json()

        } else if (response.status == 401) {

          alert("Something went wrong while authorization! Try to log in again")
          setToken()
          window.location.pathname = '/login'

        } else {

          return null

        }

      }).then(id => {
        if (id) { // If no errors

          alert("Successfully!")
          window.location.pathname = `/dispatch/${id}`

        } else { // If error appeared

          alert("Something went wrong on a server! Try again later")

        }
      })
    }
  }


  async function saveDispatch(id, token) { // Request to backend to save dispatch (only if we opened existing dispatch)
    if (title.current.value == "") { // All dispatches should have title
      alert("Write title!")
    } else {

      var dispatch = new FormData()

      dispatch.append('File', image.current.files[0])
      dispatch.append('Title', title.current.value)
      dispatch.append('Description', description.current.value)
      dispatch.append('Text', text.current.value)

      await fetch(`${API_ADRESS}/saveDispatch`, {
        method: "POST",
        headers:
        {
          "Authorization": token,
          "ID": id
        },
        body: dispatch
      }).then(response => {

        if (response.status == 202) {

          alert("Successfully!")
          window.location.reload()

        } else if (response.status == 401) {

          alert("Something went wrong while authorization! Try to log in again")
          setToken()
          window.location.pathname = '/login'

        } else {

          alert("Something went wrong on a server! Try again later")

        }

      })
    }
  }

  async function getDispatch(id) { // Request to backend to get existing disaptch

    await fetch(`${API_ADRESS}/getDispatch`, {
      headers:
      {
        "ID": `${id}`
      },
    }).then(response => {

      if (response.status == 202) {

        return response.json()

      } else {

        return null

      }

    }).then(data => {
      if (data) { // If no errors we display data

        title.current.value = data.Title
        description.current.value = data.Description
        text.current.value = data.Text
        image.current.style.backgroundImage = `url(data:${data.FileExtension};base64,${data.File})`

      } else { // If error appeared

        alert("Something went wrong on a server! Try again later")
        window.location.pathname = "/"

      }
    })
  }

  async function checkAuthor(id, token) { // Request to the backend to check User rights
    await fetch(`${API_ADRESS}/canEdit`, {
      headers:
      {
        "Authorization": token,
        "ID": id
      },
    }).then(response => {

      if (response.status == 202) {

        setCanEdit(true)

      } else if (response.status == 403) {

        setCanEdit(false)

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.pathname = '/login'

      } else {

        alert("Something went wrong on a server! Try again later")

      }

    })
  }

  async function deleteDispatch(id, token) { // Request to backend to delete dispatch
    await fetch(`${API_ADRESS}/deleteDispatch`, {
      headers:
      {
        "Authorization": token,
        "ID": id
      },
    }).then(response => {

      if (response.status == 202) {

        alert("Successfully!")
        window.location.pathname = '/'

      } else if (response.status == 403) {

        alert("You not allowed to do that!")
        window.location.reload()

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.pathname = '/login'

      } else {

        alert("Something went wrong on a server! Try again later")

      }

    })
  }

  useEffect(() => { // What will happend when page was loaded
    if (id) {
      getDispatch(id)
      if (token) {
        checkAuthor(id, token)
      }
    } else {
      setCanEdit(true)
      document.getElementById("deleteDispatch").style.display = "none";
    }
  }, [])

  useEffect(() => { // What will happend when canEdit was changed
    if (canEdit) {
      document.getElementById("isEdit").style.display = "block";
      document.getElementById("editText").style.display = "block";
    }
    else {
      document.getElementById("isEdit").style.display = "none";
      document.getElementById("editText").style.display = "none";
    }

  }, [canEdit])

  return (
    <>
      <div id='rootDispatch'>

        <div className="headerDispatch">
          <input type="checkbox" id="isEdit" className="checkSave" onClick={() => {
            if (document.getElementById("isEdit").checked) {
              document.getElementById("image1").style.pointerEvents = "all";
              document.getElementById("input1").removeAttribute("readOnly")
              document.getElementById("input2").removeAttribute("readOnly")
              document.getElementById("input3").removeAttribute("readOnly")
            } else {
              document.getElementById("image1").style.pointerEvents = "none";
              document.getElementById("input1").setAttribute("readOnly", true)
              document.getElementById("input2").setAttribute("readOnly", true)
              document.getElementById("input3").setAttribute("readOnly", true)
            }
          }} />
          <h2 className="textEdit" id="editText">Edit: </h2>
          <h3><a onClick={() => {window.location.pathname = "/"}}><u>Go back</u></a></h3>
          <input type="submit" className="Dispatch" value="Save" id="sendDispatch" onClick={() => {
            if (id) {
              saveDispatch(id, token)
            } else {
              createDispatch(token)
            }
          }} />
          <input type="submit" className="Dispatch" value="Delete" id="deleteDispatch" onClick={() => { deleteDispatch(id, token) }} />
        </div>
        <div className="mainPageDispatch">
          <input type="file" accept=".jpg,.png,.gif,.svg" onChange={
            (e) => {
              const { files } = e.target;
              if (files.length === 0) {
                return;
              }

              const file = files[0];
              const fileReader = new FileReader();

              fileReader.onload = () => {
                document.getElementById("image1").style.backgroundImage = `url(${fileReader.result})`;
              };
              fileReader.readAsDataURL(file);
            }
          } id="image1" ref={image} />
          <input type="text" placeholder="Write here your title..." className="inputTitleDispatch" maxLength="40" ref={title} id="input1" readOnly />
          <textarea className="inputDescriptionDispatch" placeholder="Write here your description..." id="input2" onKeyDown={(answ) => { if (answ.key === "Enter") { answ.preventDefault() } }} ref={description} readOnly />
          <textarea className="inputBodyDispatch" placeholder="Write here your notes..." ref={text} id="input3" readOnly />
        </div>
      </div>
    </>
  )
}

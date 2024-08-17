import { useEffect, useRef, useState } from "react"
import { API_ADRESS } from "../main.jsx"

export function ShowMiniInfo({ id }) { // Returns small dispatch info 
    const image = useRef(null)

    const [dispatchData, setdispatchData] = useState({})

    async function getDispatch(id) { // Request to backend to get dispatch ids

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
            if (data) { // If no errors

                image.current.style.backgroundImage = `url(data:${data.FileExtension};base64,${data.File})`
                setdispatchData(data)

            } else {

                alert("Something went wrong on a server! Try again later")

            }
        })
    }

    useEffect(() => {
        getDispatch(id) // Function will be called when ShowMiniInfo is called 
    }, [])


    return (
        <>
            <div id={`${id}/div`}>
                <li>
                    <img ref={image} />
                    <h1 id={id}>{dispatchData.Title}</h1>
                    <p>{dispatchData.Description}</p>
                </li>
            </div>
        </>
    )
}
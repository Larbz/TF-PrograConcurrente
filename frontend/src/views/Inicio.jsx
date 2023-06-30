import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "../context/AppProvider";
// import { TablePoints } from "./components/TablePoints";
const socket = new WebSocket("ws://localhost:3000/connection");

function Inicio() {
    const navigate = useNavigate()
    const [name, setName] = useState("");
    const [counter, setCounter] = useState(0);
    const [theSocket, setSocket] = useState(null);
    const dispatch=useDispatch();
    function startWebsocket() {
        var ws = new WebSocket("ws://localhost:3000/connection");

        ws.onmessage = function (e) {
            console.log(typeof JSON.parse(e.data))
            console.log("websocket message event:", e.data);
            if(+e.data==5){
              console.log(JSON.parse(e.data))
              console.log("COMENZAMOS!")
              navigate("/table")
            }
        };
        if(counter!=5){

            ws.onclose = function () {
                // connection closed, discard old websocket and create a new one in 5s
                ws = null;
                setTimeout(startWebsocket, 1000);
            };
        }
    }
    if(counter<5){
        startWebsocket();
    }
    const sendInfo = (e) => {
        e.preventDefault();
        socket.send(name);
        console.log(typeof JSON.parse(e.data))
        socket.onmessage = (event) => {
            console.log(JSON.parse(e.data))
            setCounter(+event.data);


        };
    };
    
    // useEffect(() => {
    //     socket.onopen = () => {
    //         console.log("Conexión WebSocket abierta");
    //     };
    //     socket.onclose = () => {
    //         console.log("disconnected");
    //     };
    // }, [counter]);
    // useEffect(() => {
    //   // Crear la instancia de WebSocket al montar el componente
    //   const newSocket = new WebSocket("ws://localhost:3000/connection");

    //   newSocket.onopen = () => {
    //     console.log("Conexión WebSocket abierta");
    //   };

    //   newSocket.onmessage = (event) => {
    //     const numero = event.data;
    //     console.log("Número recibido:", numero);
    //   };

    //   newSocket.onerror = (error) => {
    //     console.error("Error en la conexión WebSocket:", error);
    //   };

    //   // Actualizar el estado con la instancia de WebSocket
    //   setSocket(newSocket);

    //   // Cerrar la conexión al desmontar el componente
    //   return () => {
    //     newSocket.close();
    //   };
    // }, [counter]);

    // useEffect(()=>{
    //   socket.onmessage = (event) => {
    //     // const message = JSON.parse(event.data);
    //     console.log(event.data)
    //     // Verificar si el mensaje contiene el contador
    //     // if (message.counter !== undefined) {
    //       // setCounter(message.counter);

    //       // // Verificar si se alcanzó el límite de jugadores
    //       // if (message.counter === 5) {
    //       //   // Redirigir a otra vista
    //       //   history.push('/otra-vista');
    //       // }
    //     // }
    //   };

    //   return () => {
    //     socket.close();
    //   };
    // })

    return (
        <>
            <h1 style={{ textAlign: "center", fontSize: "1.6rem" }}>
                Hoop Hop Showdown – Rock Paper Scissors Hula Hoop Activity
            </h1>
            <form onSubmit={sendInfo} action="">
                <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
            </form>
            {/* <TablePoints/> */}
        </>
    );
}

export default Inicio;

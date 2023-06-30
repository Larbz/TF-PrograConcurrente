import React, { useCallback, useEffect, useRef, useState } from "react";
import { TableContainer, TablePointsContainer } from "../styles/components/TablePoints";
// import io from "socket.io-client";
// const socket = new WebSocket("ws://localhost:8050/websoc");
export const TablePoints = () => {
    const [tableInfo, setTableInfo] = useState([]);
    const [contador, setContador] = useState(0);
    // const [puntaje, setPuntaje] = useState(0);
    const fillMyTable = async () => {
        const response = await fetch("http://localhost:9000");
        const data = await response.json();
        const dataOrdered = data.sort((a, b) => b.playerPoints - a.playerPoints);
        setTableInfo(dataOrdered);
    };

    // const [message, setMessage] = useState("");
    // const [inputValue, setInputValue] = useState("");

    // useEffect(() => {
    //     const socket = io("http://localhost:8050");

    //     socket.on = () => {
    //         setMessage("Connected");
    //     };

    //     socket.onmessage = (e) => {
    //         setMessage("Get message from server: " + e.data);
    //     };
    //     // socket.on("actualizacionPuntaje",(nuevoPuntaje)=>{
    //     //     setMessage(nuevoPuntaje)
    //     // })
    //     socket.on('actualizacionPuntaje',(nuevoMessage)=>{
    //         setMessage(nuevoMessage)
    //         console.log('asd')
    //     })
    //     // socket.onclose=()=>{
    //     //     socket.close();
    //     // }
    //     // return () => {
    //     //     socket.close();
    //     // };
    // }, []);

    // const handleClick = useCallback(
    //     (e) => {
    //         e.preventDefault();

    //         socket.send(
    //             JSON.stringify({
    //                 message: inputValue,
    //             })
    //         );
    //     },
    //     [inputValue]
    // );

    // const handleChange = useCallback((e) => {
    //     setInputValue(e.target.value);
    // }, []);

    useEffect(() => {
        let intervalId;
        // const interval = setInterval(async() => {
        // fillMyTable()

        const iniciarIntervalo = () => {
            // Establecer el intervalo de 3 segundos
            intervalId = setInterval(fillMyTable, 1000);
        };

        const detenerIntervalo = () => {
            // Detener el intervalo
            clearInterval(intervalId);
        };
        fillMyTable();
        setTimeout(iniciarIntervalo, 1000);
        // }, 3000);
        return () => {
            clearTimeout();
            detenerIntervalo();
        };
    }, []);

    // const handleCambioPuntaje = () => {
    //     const newPuntaje = Math.floor(Math.random() * 100); // Generar un nuevo puntaje aleatorio
    //     socket.emit("cambio_puntaje", newPuntaje);
    // };
    return (
        <TablePointsContainer>
            <div>
                <h3>Tabla de Puntajes</h3>
            </div>
            {/* <input id="input" type="text" value={inputValue} onChange={handleChange} /> */}
            {/* <button onClick={handleClick}>Send</button> */}
            {/* <pre>{message}</pre> */}
            {tableInfo ? (
                <>
                    <TableContainer>
                        <table>
                            <thead>
                                <th></th>
                                {/* <th>Nombre</th> */}
                                <th>Posicion X</th>
                                <th>Posicion Y</th>
                                <th>Tokens ganados</th>
                                <th>Tokens perdidos</th>
                                <th>Rpc ganados</th>
                                <th>Puntos</th>
                            </thead>
                            <tbody>
                                {tableInfo.map((playerInfo, index) => (
                                    <tr key={index}>
                                        <td
                                        style={{textTransform:"uppercase"}}
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerName}
                                        </td>
                                        <td
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerPositionX}
                                        </td>
                                        <td
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerPositionY}
                                        </td>
                                        <td
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerTokensWon}
                                        </td>
                                        <td
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerTokensLost}
                                        </td>
                                        <td
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerRpcWon}
                                        </td>
                                        <td
                                            className={
                                                !playerInfo.playerPoints && "eliminated"
                                            }
                                        >
                                            {playerInfo.playerPoints}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </TableContainer>
                </>
            ) : (
                <p>Not results yet</p>
            )}
        </TablePointsContainer>
    );
};

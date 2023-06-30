import styled from "styled-components";

export const TableContainer = styled.div`
    border-collapse: collapse;
    & .eliminated {
        text-decoration: line-through;
        text-decoration-color: red;
        text-decoration-thickness: 2px;
    }
    background-color: #320b7b;
    border-radius: 30px;
    border: 1px solid #a875db;
    font-size: 1rem;
    & table {
        width: 100%;
    }
    & th {
        text-align: center;
        padding: 10px;
    }
    & td {
        text-align: center;
        padding: 10px;
    }
`;

export const TablePointsContainer = styled.div`
    width: 80%;
    margin: 30px auto;
    display: flex;
    flex-direction: column;
    gap: 20px;
    & > div:first-child {
        border: 1px solid #a875db;
        border-radius: 30px;
        background-color: #320b7b;
        padding-block: 20px;
        min-width: 450px;
        max-width: 500px;
        margin: auto;
        & h3 {
            text-align: center;
            font-size: 2rem;
            font-family: "Indie Flower", cursive;
        }
    }
`;

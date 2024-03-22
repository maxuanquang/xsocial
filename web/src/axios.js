import axios from "axios";

export const makeRequest = axios.create({
    baseURL: `${process.env.REACT_APP_API_SERVER}/api/v1`,
    withCredentials: true,
});

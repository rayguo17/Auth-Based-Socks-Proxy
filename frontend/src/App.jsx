import {useEffect, useRef, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet,ListUser,AddUser,DelUser,GetConfig} from "../wailsjs/go/app/App";
import Box from '@mui/material/Box'
import {
    Button, FilledInput,
    FormControl, FormControlLabel, FormGroup, IconButton,
    InputAdornment,
    InputLabel, List, ListItem, ListItemButton, ListItemText,
    Modal,
    OutlinedInput, Switch,
    TextField, ToggleButton, ToggleButtonGroup,
    Typography
} from "@mui/material";
import {Visibility, Delete} from "@mui/icons-material";
import {EventsOn} from "../wailsjs/runtime/runtime.js";

const dummy = [
    1,2,3,4,5,6,7,8
]
const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 600,
    bgcolor: 'background.paper',
    border: '2px solid #000',
    boxShadow: 24,
    color:"black",
    p: 4,
};

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [config,setConfig] = useState({state:"init"})
    const [name, setName] = useState('');
    const [users,setUsers] = useState([])
    const [open,setOpen] = useState(false)
    const [username,setUsername] = useState("")
    const [pwd,setPwd] = useState("")
    const [white, setWhite] = useState('black');
    const [remote,setRemote] = useState(false)
    const [remoteAddr,setRemoteAddr] = useState("")
    const [url,setUrl] = useState("")
    const [accessList,setAccessList]=useState([])
    const [pkey,setPkey] = useState("")
    const [logs,setLogs] = useState([])
    const logsTemp = useRef([])
    const logEndRef = useRef(null)
    const scrollToBottom = () => {
        logEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }
    const handleUname = (e)=>{

        setUsername(e.target.value)
    }

    useEffect(()=>{
        //initialize

        EventsOn("userChange",(data)=>{
            refreshUser()
        })
        EventsOn("log",(data)=>{
            console.log(data)
            console.log(logsTemp.current.length)
            logsTemp.current.push(data)
            setLogs(logsTemp.current)
            scrollToBottom()
        })
        EventsOn("init",(data)=>{
            console.log(data)
        })
        GetConfig().then((resp)=>{
            console.log(resp)
            setConfig({
                state: "done",
                system:resp.Data
            })
        })
        refreshUser()

    },[])
    const handleChange = (
        event,
        newAlignment,
    ) => {
        setWhite(newAlignment);
    };
    const handleAddUrl = ()=>{
        setAccessList([...accessList,url])
        setUrl("")
    }

    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);
    function handleOpen(){
        setOpen(true)
    }
    function handleClose() {
    setOpen(false)
    }
    async function refreshUser(){
        //setRefreshFlag(!refreshFlag)
        //console.log("refresh!")
        let res = await ListUser()
        //console.log(res)
        setUsers(res)
    }
    function deleteWebsite(url){
        //find and delete
        setAccessList(accessList.filter((v,i)=>{
            if (v===url){
                return false
            }
            return true
        }))
    }
    const handleDelUser =async  (delUname)=>{
        console.log(delUname)
        let res = await DelUser(delUname)
        console.log(res)
        if (res.ErrCode!==0){
            console.log("delete user fail")
            console.log(res.ErrMsg)
            return
        }else{
            refreshUser()
        }
    }
    async function handleSubmit(){
        let newUser = {
            "username":username,
            "password":pwd,
            "access":{
                "black":white!=="white",
                "black_list":white!=="white"?accessList:[],
                "white_list":white==="white"?accessList:[]
            },
            "route":{
                "type":remote?"Remote":"Direct",
                "remote":remoteAddr,
                "public_key":pkey,
                "node_id":"A868303126987902D51F2B6F06DD90038C45B119"
            }
        }
        console.log(newUser)
        let res = await AddUser(newUser)
        console.log(res)
        if (res.ErrCode !==0){
            console.log(res.ErrMsg)
        }else{
            refreshUser()
            clearForm()
        }
        //refresh

    }
    function clearForm(){
        setPkey("")
        setWhite("black")
        setRemote(false)
        setRemoteAddr("")
        setUrl("")
        setAccessList([])
        setUsername("")
        setPwd("")
        setOpen(false)
    }
    function addUserHandler(){
        //invoke modal to fill in information.
        handleOpen()
    }
    function greet() {
        Greet(name).then(updateResultText);
    }

    return (
        <div id="App" className="mx-2 h-full">
            <div>
                <h1 className="text-4xl font-bold my-2">Welcome to Yass-server</h1>
            </div>
            <div>
                <p>
                    local socks: {config.state==="init"?"":config.system.socks_port}
                </p>
                <p>
                    light server port: {config.state==="init"?"":config.system.light_config.port}
                </p>
                <p>
                    public key: {config.state==="init"?"":config.system.light_config.PublicKey}
                </p>
            </div>
            <div className="my-2">
                <button onClick={addUserHandler} className="bg-lime-600 px-1 rounded-sm mx-1">Add User</button>
                <button onClick={refreshUser} className="bg-emerald-100 text-gray-800 px-1 rounded-sm mx-1"> Refresh</button>
            </div>
            <div className="w-full overflow-auto block h-72">
                <table className="">
                    <thead>
                        <tr className="">
                            <th className="">用户名</th>
                            <th className="">上传流量(kb)</th>
                            <th className="">下载流量(kb)</th>
                            <th className="">启用</th>
                            <th className="">最后一次活动时间</th>
                            <th className="">路由</th>
                            <th className="">活跃连接</th>
                            <th className="">总连接数</th>
                            <th className="">黑白名单</th>

                        </tr>
                    </thead>
                    <tbody>
                    {users.map((v,i)=>{
                        return (
                            <tr key={i}>
                                <td>{v.Username}</td>
                                <td>{v.UploadTraffic}</td>
                                <td>{v.DownloadTraffic}</td>
                                <td>{v.Enable}</td>
                                <td>{v.LastSeen}</td>
                                <td>{v.Route}</td>
                                <td>{v.ActiveConnection}</td>
                                <td>{v.TotalConnection}</td>
                                <td>{v.Black}</td>
                                <td className="button-row"><button className="whitespace-nowrap hover:bg-cyan-400 p-1 rounded-md mx-1">查看</button></td>
                                <td onClick={()=>{handleDelUser(v.Username)}} className="button-row"><button className="whitespace-nowrap hover:bg-red-600 p-1 rounded-md mx-1">删除</button></td>
                            </tr>
                        )
                    })}
                    </tbody>
                </table>
            </div>
            <div>

                <Modal
                    open={open}
                    onClose={handleClose}
                    aria-labelledby="modal-modal-title"
                    aria-describedby="modal-modal-description"
                >
                    <Box sx={style}>
                        <h2 className="mb-2 text-3xl">User Config</h2>
                        <FormControl sx={{ mr:1,  width: '25ch' }} variant="filled">
                            <InputLabel htmlFor="outlined-adornment-uname">Username</InputLabel>
                            <FilledInput
                                id="outlined-adornment-password"
                                type='text'
                                label="Uname"
                                size="small"
                                value={username}
                                onChange={handleUname}
                            />
                        </FormControl>
                        <FormControl sx={{ mx:1,  width: '25ch' }} variant="filled">
                            <InputLabel htmlFor="outlined-adornment-password">Password</InputLabel>
                            <FilledInput
                                id="outlined-adornment-password"
                                type='password'
                                label="Password"
                                size="small"
                                value={pwd}
                                onChange={(e)=>{setPwd(e.target.value)}}
                            />
                        </FormControl>

                        <ToggleButtonGroup
                            sx={{mt:1}}
                            color="primary"
                            value={white}
                            exclusive
                            onChange={handleChange}
                            aria-label="Platform"
                        >
                            <ToggleButton value="white">WhiteList</ToggleButton>
                            <ToggleButton value="black">BlackList</ToggleButton>
                        </ToggleButtonGroup>
                        <div className='mt-1'>
                            <FormControl variant="filled">
                                <InputLabel htmlFor="outlined-adornment-url">URL</InputLabel>
                                <FilledInput
                                    id="outlined-adornment-url"
                                    type='text'
                                    size='small'
                                    value={url}
                                    onChange={(e)=>{setUrl(e.target.value)}}
                                    label="URL"/>
                            </FormControl>
                            <Button variant="text" onClick={handleAddUrl}>Add</Button>
                        </div>

                        <div className="w-1/2">
                            <List dense={true}>
                                {
                                    accessList.map((v,i)=>{
                                        return (
                                            <ListItem disablePadding key={i} secondaryAction={
                                                <IconButton onClick={()=>{deleteWebsite(v)}} edge="end" aria-label="delete">
                                                    <Delete />
                                                </IconButton>
                                            }>
                                                <ListItemText primary={v} />
                                            </ListItem>
                                        )
                                    })
                                }

                            </List>
                        </div>
                        <div>
                            <FormControlLabel control={<Switch onChange={(e)=>{setRemote(e.target.checked)}}  inputProps={{ 'aria-label': 'controlled' }}  value={remote}/>} label="Remote" />
                            {remote?(
                                <div>
                                    <FormControl variant="filled" sx={{mr:1}}>
                                        <InputLabel htmlFor="outlined-adornment-raddr">remote address</InputLabel>
                                        <FilledInput
                                            id="outlined-adornment-raddr"
                                            type='text'
                                            size='small'
                                            value={remoteAddr}
                                            onChange={(e)=>{setRemoteAddr(e.target.value)}}
                                            label="raddr"/>
                                    </FormControl>
                                    <FormControl variant="filled">
                                        <InputLabel htmlFor="outlined-adornment-pkey">public key</InputLabel>
                                        <FilledInput
                                            id="outlined-adornment-pkey"
                                            type='text'
                                            size='small'
                                            value={pkey}
                                            onChange={(e)=>{setPkey(e.target.value)}}
                                            label="pkey"/>
                                    </FormControl>
                                </div>
                            ):(<div></div>)}
                        </div>
                        <Button onClick={handleSubmit}>Submit</Button>
                    </Box>
                </Modal>
            </div>
            <div className="mt-1 " aria-label="logger">
                <h3>日志</h3>
                <div className="max-h-56 overflow-auto" >
                    {logs.map((v,i)=>{
                        return (
                            <p key={i}>{v}</p>
                        )
                    })}
                    <div ref={logEndRef}></div>
                </div>


            </div>

        </div>
    )
}

export default App

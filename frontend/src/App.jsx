import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet,ListUser} from "../wailsjs/go/app/App";
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
    const [name, setName] = useState('');
    const [users,setUsers] = useState([])
    const [open,setOpen] = useState(false)
    const [white, setWhite] = useState('black');
    const [accessList,setAccessList]=useState(["www.baidu.com","www.360.com"])
    const handleChange = (
        event,
        newAlignment,
    ) => {
        setWhite(newAlignment);
    };
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);
    function handleOpen(){
        setOpen(true)
    }
    function handleClose() {
    setOpen(false)
    }
    async function refreshUser(){
        console.log("refresh!")
        let res = await ListUser()
        console.log(res)
        setUsers(res)
    }
    function deleteWebsite(website){
        //find and delete
        console.log(website)
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
            <div className="my-2">
                <button onClick={addUserHandler} className="bg-lime-600 px-1 rounded-sm mx-1">Add User</button>
                <button onClick={refreshUser} className="bg-emerald-100 text-gray-800 px-1 rounded-sm mx-1"> Refresh</button>
            </div>
            <div className="w-full overflow-auto block h-80 ">
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
                                <td className="button-row"><button className="whitespace-nowrap hover:bg-red-600 p-1 rounded-md mx-1">删除</button></td>
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
                        <TextField id="filled-basic" label="Username" variant="filled"/>
                        <FormControl sx={{ mx:1,  width: '25ch' }} variant="filled">
                            <InputLabel htmlFor="outlined-adornment-password">Password</InputLabel>
                            <FilledInput
                                id="outlined-adornment-password"
                                type='password'
                                label="Password"/>
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

                        <Typography sx={{ mt: 1 }} variant="h6" component="div">
                            Website
                        </Typography>
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


                    </Box>
                </Modal>
            </div>

            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div>
        </div>
    )
}

export default App

import React, { useState, useEffect } from "react";

import { makeStyles } from "@material-ui/core/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";

import PowerStateButton from "./PowerStateButton.js";

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
});

function createData(name, calories, fat, carbs, protein) {
  return { name, calories, fat, carbs, protein };
}

const rows = [
  createData("Frozen yoghurt", 159, 6.0, 24, 4.0),
  createData("Ice cream sandwich", 237, 9.0, 37, 4.3),
  createData("Eclair", 262, 16.0, 24, 6.0),
  createData("Cupcake", 305, 3.7, 67, 4.3),
  createData("Gingerbread", 356, 16.0, 49, 3.9),
];

const Ec2 = (props) => {
  const classes = useStyles();

  const [data, setData] = useState([]);
  const axios = require("axios");

  useEffect(() => {
    axios
      .get("http://localhost:8000/instance?region=us-east-1")
      .then((response) => {
        console.log(response.data.instances);
        setData(response.data.instances);
      })
      .catch((err) => {
        console.log(err);
      });
  }, []);

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>InstanceId</TableCell>
            <TableCell>Name</TableCell>
            <TableCell>State</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data &&
            data.map((row) => (
              <TableRow key={row.name}>
                <TableCell scope="row">{row.Name}</TableCell>
                <TableCell>{row.Id}</TableCell>
                <TableCell>
                  <PowerStateButton
                    instanceId={row.Id}
                    powerState={row.PowerState}
                  />
                </TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default Ec2;

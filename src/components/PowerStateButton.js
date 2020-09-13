import React, { useState, useEffect } from "react";

import Button from "@material-ui/core/Button";

const PowerStateButton = (props) => {
  const [powerState, setPowerState] = useState(props.powerState);
  const [action, setAction] = useState("none");
  const axios = require("axios");

  useEffect(() => {
    if (action === "start" || action === "stop") {
      axios({
        method: "post",
        url: `http://localhost:8000/instances/${action}?region=us-east-1`,
        data: {
          instanceIds: [props.instanceId],
        },
      })
        .then((response) => {
          if (response.data === "Success") {
            setPowerState(action === "start" ? "running" : "stopped");
          }
          setAction("none");
        })
        .catch((err) => {
          setAction("none");
          console.log(err);
        });
    }
  }, [action]);

  const toggle = () => {
    if (powerState === "running") {
      setAction("stop");
    } else if (powerState === "stopped") {
      setAction("start");
    }
  };

  return (
    <Button
      onClick={toggle}
      disabled={
        (powerState !== "running" && powerState !== "stopped") ||
        action === "start" ||
        action === "stop"
      }
    >
      {powerState}
    </Button>
  );
};

export default PowerStateButton;

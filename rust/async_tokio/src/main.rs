use std::{sync::Arc, time::Duration};

use tokio::{
    sync::{
        mpsc::{self, Receiver, Sender},
        Mutex,
    },
    time::sleep,
};

async fn runCmd(tx: Sender<String>, mut rx: Receiver<String>, command: String) {

    tokio::spawn(async move {
        while let Some(msg) = rx.recv().await {
            println!("async function {}", msg);
        }
    });

    loop {
        println!("tokioroutine {}", command);
        _ = sleep(Duration::new(2, 0)).await;
        let check_output_chan = tx.send(format!("from function")).await;
        match check_output_chan {
            Ok(_) => {},
            Err(e) => eprintln!("{}",e),
        }

    }
}

#[tokio::main]
async fn main() {
    let mut hello = "test".to_string();
    let mut counter = 0;
    // I want to count the total number of running tokio threads
    let (tx_output, rx_output) = mpsc::channel::<String>(5);
    let (tx_input, rx_input) = mpsc::channel::<String>(5);

    let rx_output_arc = Arc::new(Mutex::new(rx_output));

    let _ = tokio::spawn(runCmd(tx_output, rx_input, hello.clone()));

    let rx = rx_output_arc.clone();
    let _ = tokio::spawn(async move {
        let mut guard = rx.lock().await;
        let mut counter = 0;
        loop {
            let msg = guard.recv().await;
            match msg {
                Some(msg) => {
                    if counter >= 3 {
                        println!("first tokioroutine closed");
                        drop(guard);
                        break;
                    }
                    println!("{}", msg);
                }
                _ => {}
            }
            counter += 1;
        }
    });

    // this routine will never get the messages of rx_output since the channel is locked to the
    // first routine! This routine only gets the channel if the guard from the other routine is
    // dropped (done here by drop(guard) after three iterations)

    //let rx = rx_output_arc.clone();
    //let _ = tokio::spawn(async move {
    //    let mut guard = rx.lock().await;
    //    loop {
    //        let msg = guard.recv().await;
    //        match msg {
    //            Some(msg) => println!("routine2 {}", msg),
    //            _ => {}
    //        }
    //    }
    //});

    loop {
        hello = format!("{}{}", hello, counter);
        println!("{}", hello);
        _ = tx_input.send(format!("1 via channel {}", counter)).await;
        _ = sleep(Duration::new(1, 0)).await;
        let send2 = tx_input.send(format!("2 via channel {}", counter)).await;
        match send2 {
            Ok(_) => {}
            Err(e) => eprintln!("Error while sending second message: {e}"),
        }

        
        // I can close the channel even though the runCmd function tries to send on it. In contrast
        // to go, the send command returns an error that can be gracefully handled
        let rx_output = rx_output_arc.clone();
        let mut guard = rx_output.lock().await;
        guard.close();
        println!("input routine closed");
        _ = sleep(Duration::new(7, 0)).await;
        counter += 1;
    }
}

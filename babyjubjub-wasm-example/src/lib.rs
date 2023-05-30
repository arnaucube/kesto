use babyjubjub_ark::{new_key, verify, Fq, Fr, Point, PrivateKey, Signature};

mod utils;

use wasm_bindgen::prelude::*;

// When the `wee_alloc` feature is enabled, use `wee_alloc` as the global
// allocator.
#[cfg(feature = "wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[wasm_bindgen]
extern "C" {
    fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet() {
    alert("PROVAAAAA Hello, wasm-bindings!");
}

#[wasm_bindgen]
pub fn check_eddsa_bbjj_sig() {
    let mut rng = ark_std::test_rng();
    // alert("gen new key");
    let sk = new_key(&mut rng);
    // // let sk = PrivateKey::import(
    // //     hex::decode("0001020304050607080900010203040506070809000102030405060708090001").unwrap(),
    // // )
    // // .unwrap();
    //
    let pk = sk.public();
    alert(&format!("pk: x={}, y={}", pk.x, pk.y));
    let msg = Fq::from(5_u32);
    alert(&format!("msg: {:?}", msg.to_string()));
    let sig = sk.sign(msg.clone()).unwrap();
    alert(&format!(
        "signature:\ns:{}, r.x: {}, r.y: {}",
        sig.s, sig.r_b8.x, sig.r_b8.x
    ));
    let v = verify(pk.clone(), sig.clone(), msg.clone());
    assert_eq!(v, true);
    alert(&format!("signature verification: {}", v));
}

#[test]
fn test_compat() {
    let mut rng = ark_std::test_rng();
    let sk = new_key(&mut rng);
    // // let sk = PrivateKey::import(
    // //     hex::decode("0001020304050607080900010203040506070809000102030405060708090001").unwrap(),
    // // )
    // // .unwrap();
    //
    let pk = sk.public();
    println!("pk: x={}, y={}", pk.x, pk.y);
    let msg = Fq::from(5_u32);
    println!("msg: {:?}", msg.to_string());
    let sig = sk.sign(msg.clone()).unwrap();
    println!(
        "signature:\ns:{}, r.x: {}, r.y: {}",
        sig.s, sig.r_b8.x, sig.r_b8.x
    );
    let v = verify(pk.clone(), sig.clone(), msg.clone());
    assert_eq!(v, true);
}

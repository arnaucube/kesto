use criterion::{criterion_group, criterion_main, Criterion};

// JubJub
use dusk_bls12_381::Scalar;
use eddsa::{KeyPair, Message, PublicKey};
extern crate rand7;

// BabyJubJub
extern crate rand;
#[macro_use]
extern crate ff;
use ff::*;
extern crate num;
extern crate num_bigint;
use babyjubjub_rs::{utils, Point};
use num_bigint::{BigInt, Sign, ToBigInt};

fn criterion_benchmark(c: &mut Criterion) {
    let mut m: [u8; 32] = rand::random::<[u8; 32]>();
    m[31] = 0;
    println!("m {:?}", m);

    // JubJub
    let keypair = KeyPair::new(&mut rand7::thread_rng()).unwrap();
    let message = Message(Scalar::from_bytes(&m).unwrap());
    c.bench_function("JubJub EdDSA sign", |b| b.iter(|| keypair.sign(&message)));
    let a = keypair.sign(&message);
    c.bench_function("JubJub EdDSA verify", |b| {
        b.iter(|| a.verify(&message, &keypair.public_key))
    });

    // BabyJubJub
    let sk = babyjubjub_rs::new_key();
    let pk = sk.public().unwrap();
    let msg = BigInt::from_bytes_le(Sign::Plus, &m);
    c.bench_function("BabyJubJub EdDSA sign", |b| b.iter(|| sk.sign(msg.clone())));
    let sig = sk.sign(msg.clone()).unwrap();
    c.bench_function("BabyJubJub EdDSA verify", |b| {
        b.iter(|| babyjubjub_rs::verify(pk.clone(), sig.clone(), msg.clone()))
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);

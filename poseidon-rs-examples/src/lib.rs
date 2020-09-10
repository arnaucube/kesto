#[cfg(test)]
mod tests {
    use ff::PrimeField;
    use poseidon_rs::{Fr, Poseidon};

    // for parsing hex
    extern crate num;
    extern crate num_bigint;
    use num_bigint::BigInt;

    #[test]
    fn test_usage() {
        let v: Fr = Fr::from_str(
            "11043376183861534927536506085090418075369306574649619885724436265926427398571",
        )
        .unwrap();
        let mut to_hash: Vec<Fr> = Vec::new();
        to_hash.push(v);

        let poseidon = Poseidon::new();
        let h = poseidon.hash(to_hash).unwrap();
        assert_eq!(
            h.to_string(),
            "Fr(0x28410c403c92a9f18d1f27b22218b3649b3be8640dc160ad53bd21cf02f98d81)"
        );
    }

    #[test]
    fn test_usage_hex() {
        let b: BigInt = BigInt::parse_bytes(
            b"186a5454a7c47c73dfc74ac32ea40a57d27eeb4e2bfc6551dd7b66686d3fd1ab", // same value than in previous test, but in hex
            16,
        )
        .unwrap();

        let v: Fr = Fr::from_str(&b.to_string()).unwrap();
        let mut to_hash: Vec<Fr> = Vec::new();
        to_hash.push(v);

        let poseidon = Poseidon::new();
        let h = poseidon.hash(to_hash).unwrap();
        assert_eq!(
            h.to_string(),
            "Fr(0x28410c403c92a9f18d1f27b22218b3649b3be8640dc160ad53bd21cf02f98d81)"
        );
    }

    #[test]
    fn test_usage_bytes() {
        let msg = "hello";
        let b: BigInt = BigInt::parse_bytes(msg.as_bytes(), 10).unwrap();
        let v: Fr = Fr::from_str(&b.to_string()).unwrap();

        let mut to_hash: Vec<Fr> = Vec::new();
        to_hash.push(v);

        let poseidon = Poseidon::new();
        let h = poseidon.hash(to_hash).unwrap();
        assert_eq!(
            h.to_string(),
            "Fr(0x28410c403c92a9f18d1f27b22218b3649b3be8640dc160ad53bd21cf02f98d81)"
        );
    }
}

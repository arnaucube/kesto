#[cfg(test)]
mod tests {
    use poseidon_rs::Poseidon;
    use rustc_hex::ToHex;

    #[test]
    fn test_output_size() {
        let poseidon = Poseidon::new();
        let msg = "45";
        let h = poseidon.hash_bytes(msg.as_bytes().to_vec()).unwrap();
        println!("bigint {:?}", h.to_string());
        println!("length {:?}", h.to_bytes_be().1.len());
        println!("bytes {:?}", h.to_bytes_be().1);
        assert_eq!(h.to_bytes_be().1.len(), 31);
        println!("hex {:?}", h.to_bytes_be().1.to_hex());
    }
}

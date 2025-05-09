import styles from './input.module.css'

const Input = ({ className, error, label, type = "text", value, onChange, placeholder, filePreview, ...props }) => {
    return (
        <div className={styles.inputWrapper}>
            {label && <label className={styles.label}>{label}</label>}
            <input
                className={`${styles.input} ${type==="file" ? styles.inputFile : ""} ${className || ""}`}
                type={type}
                onChange={onChange}
                placeholder={placeholder}
                {...(type !== "file" ? { value: value } : {})}
                {...props}
            />
            {type === "file" && filePreview && (
                <span className={styles.filePreview}>{filePreview}</span>
            )}
            {error && <span className={styles.error}>{error}</span>}
        </div>
    );
};

export { Input };
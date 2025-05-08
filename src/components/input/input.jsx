import styles from './input.module.css'

const Input = ({ error, label, type = "text", value, onChange, placeholder, ...props }) => {
    return (
        <div className={styles.inputWrapper}>
            {label && <label className={styles.label}>{label}</label>}
            <input
                className={styles.input}
                type={type}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                {...props}
            />
            {error && <span>{error}</span>}
        </div>
    );
};

export { Input };
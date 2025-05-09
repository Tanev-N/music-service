import styles from "./button.module.css";

const Button = ({ text, onClick, type, size = "medium" }) => {
  return (
    <button
      className={`${styles.button} ${styles[size]} ${
        type === "submit" ? styles["sign-up"] : ""
      } ${type === "delete" ? styles["delete"] : ""}`}
      onClick={onClick}
    >
      {text}
    </button>
  );
};

export { Button };

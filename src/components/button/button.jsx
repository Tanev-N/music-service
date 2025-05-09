import styles from "./button.module.css";

const Button = ({ text, onClick, type }) => {
  return (
    <button
      className={`${styles.button} ${
        type === "submit" ? styles["sign-up"] : ""
      } ${type === "delete" ? styles["delete"] : ""}`}
      onClick={onClick}
    >
      {text}
    </button>
  );
};

export { Button };

import styles from "./button.module.css";

const Button = ({ text, onClick, type }) => {
  return (
    <button
      className={`${styles.button} ${
        type === "submit" ? styles["sign-up"] : ""
      }`}
      onClick={onClick}
    >
      {text}
    </button>
  );
};

export { Button };

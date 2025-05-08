import styles from "./button.module.css";

const Button = ({ text, onClick }) => {
  return <button className={styles.button} onClick={onClick}>{text}</button>;
};

export {Button}

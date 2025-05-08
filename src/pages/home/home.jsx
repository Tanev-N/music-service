import styles from "./home.module.css";

const HomePage = () => {
  return (
    <main className={styles.main}>
      <div className={styles.mainRow}>
        <section className={styles.promoContent}>
          <h1>Музыка и книги всегда с вами!</h1>
          <p className={styles.subtitle}>
            Слушайте любимые треки и аудиокниги без ограничений. Создавайте свои коллекции и наслаждайтесь отличным звуком в любое время.
          </p>
        </section>
        <img
        //   src="/music-icon.png"
          src="/note.png"
          alt="Наушники"
          className={styles.promoImage}
        />
      </div>
      <div className={styles.mainRow}>
        <img
          src="/peorple-women.svg"
          alt="Наушники"
          className={styles.promoImage}
        />
        <section className={styles.promoContent}>
          <h1>Ваша музыка — ваш стиль жизни</h1>
          <p className={styles.subtitle}>
            Откройте для себя новые жанры, делитесь плейлистами с друзьями и слушайте без рекламы и интернета.
          </p>
        </section>
      </div>
      <div className={styles.mainRow}>
        <section className={styles.promoContent}>
          <h1>Погружайтесь в мир звука</h1>
          <p className={styles.subtitle}>
            Тысячи исполнителей и авторов ждут вас. Сохраняйте избранное и открывайте новое каждый день.
          </p>
        </section>
        <img
          src="/green-man.svg"
          alt="Наушники"
          className={styles.promoImage}
        />
      </div>
      
    </main>
  );
};

export default HomePage;

import styles from './buttons.module.scss';

const mainButtonClass = [styles.mmPrimaryButton];

export const Basic = ({ children, mainClassStyles: mainStyles = [mainButtonClass], className, ...props }) => {
  const mainClasses = Array.isArray(mainStyles) ? mainStyles : [mainStyles];
  const classNameValue = [mainClasses, className || []].flat().join(' ');
  return (
    <button {...props} className={classNameValue}>
      {children}
    </button>
  );
};

Basic.WithNavAction = ({ href, ...props }) => {
  return <Basic onClick={() => (window.location.href = href)} {...props} />;
};

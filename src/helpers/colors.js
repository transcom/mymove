import colors from '../shared/styles/colors.scss';

const getColor = (color) => {
  const colorKeys = Object.keys(colors);
  if (colorKeys.includes(color)) {
    // The use of colors[color] triggers a security warning from our eslint security plugin.
    // However, since we verify inputs against imported colors and this function is not used where
    // users input color we are diabling the warning.
    // eslint-disable-next-line security/detect-object-injection
    return colors[color];
  }
  /* eslint-disable no-console */
  console.error(`Could not find the color: "${color}".\n\nAvailable colors:\n\n${colorKeys.join('\n')}`);
  /* eslint-enable no-console */
  return colors.base;
};

export default getColor;

// eslint-disable-next-line import/no-extraneous-dependencies
import 'happo-plugin-storybook/register';

import './storybook.scss';
import '../src/index.scss';
import '../src/ghc_index.scss';

import '../src/icons';

export const parameters = {
  options: {
    storySort: {
      order: ['Global', 'Components', 'Office Components', 'Customer Components', 'Customer Pages', 'Samples', 'Scenes',],
    },
  },
};

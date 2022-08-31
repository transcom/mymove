const { v4 } = require('uuid');

const createBaseProGear = () => {
  return {
    id: v4(),
  };
};

export default createBaseProGear;

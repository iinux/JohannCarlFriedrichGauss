import seaborn as sns
import matplotlib.pyplot as plt

# %matplotlib inline
titanic = sns.load_dataset('titanic')
sns.barplot(x='class', y='survived', data=titanic)

if __name__ == '__main__':
    pass

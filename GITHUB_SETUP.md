# Push to GitHub and open PR

1. **Create the repository on GitHub**
   - Go to https://github.com/new
   - Repository name: **grpc-go-fx**
   - Public, no README / .gitignore / license (we already have them)
   - Create repository

2. **Add remote and push** (replace `felipelimaa` with your GitHub username)

   ```bash
   cd /Users/felipe.lima/Developer/grpc-go-fx
   git remote add origin https://github.com/felipelimaa/grpc-go-fx.git
   git push -u origin feat-initial-setup
   git push origin main
   ```

   Or with SSH:

   ```bash
   git remote add origin git@github.com:felipelimaa/grpc-go-fx.git
   git push -u origin feat-initial-setup
   git push origin main
   ```

3. **Open a Pull Request**
   - On GitHub: **Compare & pull request** for branch `feat-initial-setup`
   - Base: **main** ← Compare: **feat-initial-setup**
   - Title: `feat: initial setup`
   - Create pull request
